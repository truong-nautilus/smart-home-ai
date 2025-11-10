package openai

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/truong-nautilus/smart-home-ai/internal/domain"
)

// KhachHang triển khai các dịch vụ dựa trên OpenAI
type KhachHang struct {
	client *openai.Client
}

// NewKhachHang tạo một OpenAI client mới
func NewKhachHang(apiKey string) *KhachHang {
	return &KhachHang{
		client: openai.NewClient(apiKey),
	}
}

// ChuyenGloi triển khai BoNhanDienGiongNoi
func (c *KhachHang) ChuyenGloi(ctx context.Context, duongDanAmThanh string) (*domain.BanChuyenGloi, error) {
	file, err := os.Open(duongDanAmThanh)
	if err != nil {
		return nil, fmt.Errorf("không thể mở file âm thanh: %w", err)
	}
	defer file.Close()

	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: file.Name(),
		Reader:   file,
	}

	resp, err := c.client.CreateTranscription(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("không thể tạo bản chuyển glọi: %w", err)
	}

	return &domain.BanChuyenGloi{
		VanBan:    resp.Text,
		NgonNgu:   resp.Language,
		ThoiLuong: float64(resp.Duration),
		ThoiGian:  time.Now(),
	}, nil
}

// PhanTichDaPhuongThuc triển khai TroLyAI
func (c *KhachHang) PhanTichDaPhuongThuc(ctx context.Context, vanBan string, duongDanAnh string) (*domain.PhanHoiAI, error) {
	// Đọc và mã hóa ảnh sang base64
	duLieuAnh, err := os.ReadFile(duongDanAnh)
	if err != nil {
		return nil, fmt.Errorf("không thể đọc file ảnh: %w", err)
	}
	anhBase64 := base64.StdEncoding.EncodeToString(duLieuAnh)

	req := openai.ChatCompletionRequest{
		Model: "gpt-4o-mini",
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleUser,
				MultiContent: []openai.ChatMessagePart{
					{
						Type: openai.ChatMessagePartTypeText,
						Text: fmt.Sprintf("Tôi đã nói: \"%s\"\n\nHãy mô tả những gì bạn thấy trong hình ảnh và trả lời câu hỏi của tôi dựa trên những gì tôi đã nói.", vanBan),
					},
					{
						Type: openai.ChatMessagePartTypeImageURL,
						ImageURL: &openai.ChatMessageImageURL{
							URL:    fmt.Sprintf("data:image/jpeg;base64,%s", anhBase64),
							Detail: openai.ImageURLDetailAuto,
						},
					},
				},
			},
		},
		MaxTokens: 500,
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("không thể tạo chat completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("không nhận được phản hồi từ GPT")
	}

	return &domain.PhanHoiAI{
		VanBan:   resp.Choices[0].Message.Content,
		MoHinh:   resp.Model,
		ThoiGian: time.Now(),
	}, nil
}

// TongHop triển khai BoTongHopGiongNoi
func (c *KhachHang) TongHop(ctx context.Context, vanBan string, duongDanDauRa string) (*domain.AmThanhTongHop, error) {
	req := openai.CreateSpeechRequest{
		Model:          openai.TTSModel1,
		Input:          vanBan,
		Voice:          openai.VoiceAlloy,
		ResponseFormat: openai.SpeechResponseFormatMp3,
	}

	resp, err := c.client.CreateSpeech(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("không thể tạo giọng nói: %w", err)
	}
	defer resp.Close()

	tepDauRa, err := os.Create(duongDanDauRa)
	if err != nil {
		return nil, fmt.Errorf("không thể tạo file đầu ra: %w", err)
	}
	defer tepDauRa.Close()

	if _, err = io.Copy(tepDauRa, resp); err != nil {
		return nil, fmt.Errorf("không thể ghi dữ liệu âm thanh: %w", err)
	}

	return &domain.AmThanhTongHop{
		DuongDanAmThanh: duongDanDauRa,
		VanBan:          vanBan,
		GiongNoi:        string(openai.VoiceAlloy),
		ThoiGian:        time.Now(),
	}, nil
}
