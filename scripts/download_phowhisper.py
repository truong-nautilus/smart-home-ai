#!/Users/phamthetruong/phowhisper-env/bin/python3
"""
Pre-download PhoWhisper model to cache with retry and timeout handling
"""
import os
print("🔄 Đang pre-download PhoWhisper model...")

# Tăng timeout và retry
os.environ['HF_HUB_DOWNLOAD_TIMEOUT'] = '300'  # 5 phút timeout

from transformers import WhisperForConditionalGeneration, WhisperProcessor
from huggingface_hub import snapshot_download

model_id = "vinai/PhoWhisper-small"

print(f"📥 Downloading model: {model_id}")
print("⏳ Đây có thể mất vài phút (timeout: 5 phút)...")
print("🌐 Kết nối đến HuggingFace...")

try:
    # Thử download toàn bộ repository trước
    print("\n📦 Downloading all files...")
    snapshot_download(
        repo_id=model_id,
        resume_download=True,
        local_files_only=False,
    )
    print("✅ All files downloaded")
    
    # Load processor và model để verify
    print("\n1/2 Loading processor...")
    processor = WhisperProcessor.from_pretrained(model_id)
    print("✅ Processor loaded")
    
    print("\n2/2 Loading model...")
    model = WhisperForConditionalGeneration.from_pretrained(model_id)
    print("✅ Model loaded")
    
    print("\n🎉 PhoWhisper model đã được download và cache thành công!")
    print("📂 Location: ~/.cache/huggingface/hub/")
    
except Exception as e:
    print(f"\n❌ Lỗi: {e}")
    print("\n💡 Gợi ý:")
    print("   1. Kiểm tra kết nối internet")
    print("   2. Thử lại sau vài phút")
    print("   3. Hoặc dùng VPN nếu kết nối chậm")
    exit(1)
