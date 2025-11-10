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
os.environ['TOKENIZERS_PARALLELISM'] = 'false'  # Tránh warning
import warnings
warnings.filterwarnings('ignore')

# Cache pipeline globally để không phải load lại
_pipeline_cache = None

def get_pipeline():
    """Get or create cached pipeline"""
    global _pipeline_cache
    if _pipeline_cache is None:
        from transformers import pipeline
        import torch
        
        device = "cuda:0" if torch.cuda.is_available() else "cpu"
        torch_dtype = torch.float16 if torch.cuda.is_available() else torch.float32
        model_id = os.getenv("PHOWHISPER_MODEL", "vinai/PhoWhisper-small")
        
        _pipeline_cache = pipeline(
            "automatic-speech-recognition",
            model=model_id,
            torch_dtype=torch_dtype,
            device=device,
        )
    return _pipeline_cache

def transcribe_audio(audio_path: str) -> str:
    """Transcribe audio file using PhoWhisper"""
    try:
        # Sử dụng cached pipeline
        pipe = get_pipeline()
        
        # Transcribe với tối ưu tốc độ
        # Bỏ chunk_length_s cho audio ngắn để xử lý nhanh hơn
        result = pipe(
            audio_path,
            return_timestamps=False,
            generate_kwargs={
                "language": "vietnamese", 
                "task": "transcribe",
                "num_beams": 1,              # Greedy decoding - nhanh nhất
                "do_sample": False,           # Không sampling - deterministic
                "max_length": 448,            # Giới hạn độ dài output
            },
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
