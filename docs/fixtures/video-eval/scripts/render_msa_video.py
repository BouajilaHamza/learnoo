import subprocess
import os

SLIDES = [1, 2, 3, 4, 5]
TEMP_DIR = "/tmp/learnoo_video"
MSA_DIR = os.path.join(TEMP_DIR, "msa")

def get_duration(audio_file):
    cmd = [
        "ffprobe", "-v", "error", "-show_entries", "format=duration",
        "-of", "default=noprint_wrappers=1:nokey=1", audio_file
    ]
    res = subprocess.run(cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True, check=True)
    return float(res.stdout.strip())

# 1. Concat audio
print("Concatenating MSA audio...")
inputs = []
for s in SLIDES:
    inputs.extend(["-i", os.path.join(MSA_DIR, f"slide{s}_msa.mp3")])

concat_filter = "".join([f"[{i}:a]" for i in range(len(SLIDES))]) + f"concat=n={len(SLIDES)}:v=0:a=1[outa]"
full_audio = "/home/hamza/Desktop/cours_nombres_complexes_msa.mp3"
subprocess.run([
    "ffmpeg", "-y", *inputs,
    "-filter_complex", concat_filter,
    "-map", "[outa]", full_audio
], check=True)
print(f"✓ Saved master audio to {full_audio}")

# 2. Render each video segment
segment_files = []
for s in SLIDES:
    png = os.path.join(TEMP_DIR, f"slide{s}.png")
    audio = os.path.join(MSA_DIR, f"slide{s}_msa.mp3")
    dur = get_duration(audio) + 0.6  # padding
    seg_mp4 = os.path.join(MSA_DIR, f"slide{s}_msa.mp4")
    
    cmd = [
        "ffmpeg", "-y",
        "-loop", "1", "-i", png,
        "-i", audio,
        "-c:v", "libx264", "-tune", "stillimage", "-c:a", "aac", "-b:a", "192k",
        "-pix_fmt", "yuv420p",
        "-t", str(dur),
        seg_mp4
    ]
    subprocess.run(cmd, check=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    segment_files.append(seg_mp4)
    print(f"✓ Rendered segment {s}")

# 3. Concat master video
concat_txt = os.path.join(MSA_DIR, "concat.txt")
with open(concat_txt, "w") as f:
    for seg in segment_files:
        f.write(f"file '{seg}'\n")

master_mp4 = "/home/hamza/Desktop/cours_nombres_complexes_msa_video.mp4"
subprocess.run([
    "ffmpeg", "-y", "-f", "concat", "-safe", "0",
    "-i", concat_txt,
    "-c", "copy",
    master_mp4
], check=True)
print(f"✓ Saved master video to {master_mp4}")
