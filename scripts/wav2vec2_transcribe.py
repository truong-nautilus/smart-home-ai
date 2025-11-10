#!/Users/phamthetruong/phowhisper-env/bin/python3
"""
Wav2Vec2 transcription script - FAST Vietnamese ASR
Model: Wav2Vec2-Base-Vietnamese-160h by NTT123
"""
import sys
import os

# Tắt warnings
os.environ['TF_CPP_MIN_LOG_LEVEL'] = '3'
import warnings
warnings.filterwarnings('ignore')

def transcribe_audio(audio_path: str) -> str:
    """Transcribe audio file using Wav2Vec2-Base-Vietnamese-160h"""
    try:
        from transformers import Wav2Vec2ForCTC, Wav2Vec2Processor
        import torch
        import librosa
        
        # Load model (nhẹ và nhanh hơn PhoWhisper rất nhiều)
        model_id = "nguyenvulebinh/wav2vec2-base-vietnamese-250h"
        processor = Wav2Vec2Processor.from_pretrained(model_id)
        model = Wav2Vec2ForCTC.from_pretrained(model_id)
        
        # Load audio
        speech, rate = librosa.load(audio_path, sr=16000)
        
        # Transcribe
        inputs = processor(speech, sampling_rate=16000, return_tensors="pt", padding=True)
        
        with torch.no_grad():
            logits = model(inputs.input_values).logits
        
        predicted_ids = torch.argmax(logits, dim=-1)
        transcription = processor.batch_decode(predicted_ids)[0]
        
        return transcription.strip()
        
    except Exception as e:
        print(f"Wav2Vec2 Error: {e}", file=sys.stderr)
        return ""

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: wav2vec2_transcribe.py <audio_file>", file=sys.stderr)
        sys.exit(1)
    
    audio_file = sys.argv[1]
    
    text = transcribe_audio(audio_file)
    if text:
        print(text)
    else:
        sys.exit(1)
