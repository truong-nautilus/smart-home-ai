#!/Users/phamthetruong/phowhisper-env/bin/python3
"""
PhoWhisper transcription script - Advanced Vietnamese ASR
Model: vinai/PhoWhisper-small (157M params)
Optimized for Vietnamese language with high accuracy
"""
import sys
import os

# Tắt tất cả warnings và logging
os.environ['TF_CPP_MIN_LOG_LEVEL'] = '3'
os.environ['TRANSFORMERS_VERBOSITY'] = 'error'
import warnings
warnings.filterwarnings('ignore')

def transcribe_audio(audio_path: str) -> str:
    """Transcribe audio file using PhoWhisper"""
    try:
        from transformers import pipeline
        import torch
        
        # Kiểm tra GPU/CPU
        device = "cuda:0" if torch.cuda.is_available() else "cpu"
        torch_dtype = torch.float16 if torch.cuda.is_available() else torch.float32
        
        # Load PhoWhisper model
        # Model này được tối ưu cho tiếng Việt, cho kết quả tốt hơn Whisper gốc
        model_id = os.getenv("PHOWHISPER_MODEL", "vinai/PhoWhisper-small")
        
        pipe = pipeline(
            "automatic-speech-recognition",
            model=model_id,
            dtype=torch_dtype,  # Sử dụng dtype thay vì torch_dtype
            device=device,
        )
        
        # Transcribe
        result = pipe(
            audio_path,
            chunk_length_s=30,  # Xử lý từng đoạn 30 giây
            batch_size=8,       # Batch size tùy theo RAM/VRAM
            return_timestamps=False,
            generate_kwargs={"language": "vietnamese", "task": "transcribe"},
        )
        
        transcription = result["text"].strip()
        return transcription
        
    except Exception as e:
        print(f"PhoWhisper Error: {e}", file=sys.stderr)
        return ""

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: phowhisper_transcribe.py <audio_file>", file=sys.stderr)
        sys.exit(1)
    
    audio_file = sys.argv[1]
    
    if not os.path.exists(audio_file):
        print(f"Error: File not found: {audio_file}", file=sys.stderr)
        sys.exit(1)
    
    text = transcribe_audio(audio_file)
    if text:
        # Chỉ print kết quả ra stdout, không có gì khác
        print(text, flush=True)
    else:
        sys.exit(1)
