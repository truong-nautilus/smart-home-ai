#!/usr/bin/env python3
"""
Gesture detector - Nhận diện cử chỉ 2 ngón tay từ camera (chạy ẩn, không hiển thị cửa sổ)
Sử dụng MediaPipe Hands để detect số ngón tay
"""
import cv2
import mediapipe as mp
import sys
import time
import os

def count_fingers(hand_landmarks):
    """Đếm số ngón tay đang giơ lên"""
    # Tip của các ngón tay (không tính cổ tay)
    finger_tips = [8, 12, 16, 20]  # index, middle, ring, pinky
    thumb_tip = 4
    
    count = 0
    
    # Kiểm tra ngón cái (so sánh x coordinate vì ngón cái nằm ngang)
    if hand_landmarks.landmark[thumb_tip].x < hand_landmarks.landmark[thumb_tip - 1].x:
        count += 1
    
    # Kiểm tra các ngón còn lại (so sánh y coordinate)
    for tip in finger_tips:
        if hand_landmarks.landmark[tip].y < hand_landmarks.landmark[tip - 2].y:
            count += 1
    
    return count

def main():
    # Khởi tạo MediaPipe Hands
    mp_hands = mp.solutions.hands
    hands = mp_hands.Hands(
        static_image_mode=False,
        max_num_hands=1,
        min_detection_confidence=0.5,
        min_tracking_confidence=0.5
    )
    
    # Mở camera
    print("Đang mở camera...", file=sys.stderr)
    cap = cv2.VideoCapture(0)
    
    if not cap.isOpened():
        print("ERROR: Cannot open camera", file=sys.stderr)
        sys.exit(1)
    
    print("Waiting for 2 fingers gesture...", file=sys.stderr, flush=True)
    
    detected = False
    frame_count = 0
    
    # Luôn hiển thị cửa sổ để dễ test
    show_window = True
    window_name = "Gesture Detection - Show 2 fingers"
    cv2.namedWindow(window_name)
    
    while True:
        ret, frame = cap.read()
        if not ret:
            print("ERROR: Failed to read frame from camera", file=sys.stderr, flush=True)
            time.sleep(0.5)
            continue
        
        frame_count += 1
        
        # Chuyển BGR sang RGB
        rgb_frame = cv2.cvtColor(frame, cv2.COLOR_BGR2RGB)
        results = hands.process(rgb_frame)
        
        finger_count = 0
        
        if results.multi_hand_landmarks:
            for hand_landmarks in results.multi_hand_landmarks:
                finger_count = count_fingers(hand_landmarks)
                
                # Log occasionally để biết script đang chạy
                if frame_count % 30 == 0:  # Log mỗi 1 giây (30 frames)
                    print(f"Detected {finger_count} fingers", file=sys.stderr, flush=True)
                
                if finger_count == 2:
                    print("DETECTED_2_FINGERS", file=sys.stdout, flush=True)
                    detected = True
                    break
        
        # Hiển thị frame với thông tin
        cv2.putText(frame, "Show 2 fingers to start", (10, 30), 
                    cv2.FONT_HERSHEY_SIMPLEX, 0.7, (0, 255, 0), 2)
        cv2.putText(frame, f"Fingers: {finger_count}", (10, 70), 
                    cv2.FONT_HERSHEY_SIMPLEX, 0.7, (255, 255, 0), 2)
        
        if finger_count == 2:
            cv2.putText(frame, "2 FINGERS DETECTED!", (10, 110), 
                        cv2.FONT_HERSHEY_SIMPLEX, 0.7, (0, 255, 0), 2)
        
        cv2.imshow(window_name, frame)
        
        if cv2.waitKey(1) & 0xFF == ord('q'):
            break
        
        if detected:
            # Hiển thị kết quả 0.5 giây trước khi đóng
            time.sleep(0.5)
            break
        
        # Giảm CPU usage - chỉ process 20 fps thay vì 30 fps
        time.sleep(0.05)
    
    cv2.destroyAllWindows()
    cap.release()
    hands.close()
    
    sys.exit(0 if detected else 1)

if __name__ == "__main__":
    main()
