import os
import subprocess
import json

SCRIPT_DIR = "/tmp/learnoo_video"
os.makedirs(SCRIPT_DIR, exist_ok=True)

slides_data = [
    {
        "id": "slide1",
        "title": "Nombres Complexes · الأعداد المركبة",
        "subtitle": "Forme algébrique, conjugué et module · مراجعة شاملة للباكالوريا",
        "badge": "BAC TUNISIE 2026 · باكالوريا تونس",
        "content_html": """
        <div style="margin-top: 40px; background: rgba(2, 132, 199, 0.08); border: 2px solid rgba(2, 132, 199, 0.2); border-radius: 24px; padding: 36px 48px; display: inline-block;">
            <div style="font-size: 38px; color: #0284c7; font-weight: 800; font-family: 'Cairo', sans-serif; direction: rtl;">
                شرح مبسّط بالدارجة التونسية للمفاهيم الأساسية
            </div>
            <div style="font-size: 24px; color: #475569; margin-top: 14px; font-family: 'Plus Jakarta Sans', sans-serif;">
                Section Mathématiques & Sciences Expérimentales
            </div>
        </div>
        """,
        "voice_text": "عسلامة ومرحباً بيكم تلامذة الباكالوريا في تونس. في هذا الفيديو القصير باش نفهمو مع بعضنا درس الأعداد المركبة، الشكل الجبري والمعيار، بالدارجة التونسية وبأبسط طريقة."
    },
    {
        "id": "slide2",
        "title": "1. La Forme Algébrique · الشكل الجبري",
        "subtitle": "تعريف العدد المركب والقاعدة الذهبية",
        "badge": "DÉFINITION FONDAMENTALE",
        "content_html": """
        <div style="display: flex; gap: 40px; justify-content: center; margin-top: 30px;">
            <div style="flex: 1; background: #ffffff; border: 2px solid #e2e8f0; border-radius: 20px; padding: 30px; box-shadow: 0 10px 25px rgba(0,0,0,0.05);">
                <div style="font-size: 24px; color: #0f172a; font-weight: 700; margin-bottom: 15px;">La Règle d'or · القاعدة الأساسية</div>
                <div style="font-size: 42px; color: #dc2626; font-weight: 800; font-family: monospace; background: #fef2f2; padding: 15px; border-radius: 12px;">
                    i² = -1
                </div>
                <div style="font-size: 20px; color: #64748b; margin-top: 15px; direction: rtl; font-family: 'Cairo', sans-serif;">
                    عدد جديد يسمح بحل المعادلات التي ليس لها حل في ℝ.
                </div>
            </div>
            <div style="flex: 1; background: #ffffff; border: 2px solid #e2e8f0; border-radius: 20px; padding: 30px; box-shadow: 0 10px 25px rgba(0,0,0,0.05);">
                <div style="font-size: 24px; color: #0f172a; font-weight: 700; margin-bottom: 15px;">Forme Algébrique · الشكل الجبري</div>
                <div style="font-size: 42px; color: #0284c7; font-weight: 800; font-family: monospace; background: #f0f9ff; padding: 15px; border-radius: 12px;">
                    z = a + i b
                </div>
                <div style="font-size: 20px; color: #334155; margin-top: 15px; direction: rtl; font-family: 'Cairo', sans-serif; line-height: 1.6;">
                    • <b>a</b> هو الجزء الحقيقي (Partie réelle): Re(z) = a<br>
                    • <b>b</b> هو الجزء التخيلي (Partie imaginaire): Im(z) = b
                </div>
            </div>
        </div>
        """,
        "voice_text": "نبداو بالقاعدة الأساسية. في مجموعة الأعداد المركبة، فما عدد جديد اسمو إي، بحيث إي قوّة اثنين تساوي إلا واحد. وأي عدد مركب زاد يتكتب على الشكل الجبري: زاد يساوي أ مع إي ب. أ هو الجزء الحقيقي، وب هو الجزء التخيلي."
    },
    {
        "id": "slide3",
        "title": "2. Conjugué et Module · المرافق والمعيار",
        "subtitle": "المفاهيم الهندسية والمسافة في المستوى المركب",
        "badge": "GÉOMÉTRIE DANS LE PLAN COMPLEXE",
        "content_html": """
        <div style="display: flex; gap: 40px; justify-content: center; margin-top: 30px;">
            <div style="flex: 1; background: #ffffff; border: 2px solid #e2e8f0; border-radius: 20px; padding: 30px; box-shadow: 0 10px 25px rgba(0,0,0,0.05);">
                <div style="font-size: 24px; color: #0f172a; font-weight: 700; margin-bottom: 15px;">Le Conjugué · مرافق عدد مركب</div>
                <div style="font-size: 38px; color: #0284c7; font-weight: 800; font-family: monospace; background: #f0f9ff; padding: 15px; border-radius: 12px;">
                    z̄ = a - i b
                </div>
                <div style="font-size: 20px; color: #475569; margin-top: 15px; direction: rtl; font-family: 'Cairo', sans-serif;">
                    نبدلو علامة الجزء التخيلي فقط. هندسياً هو تناظر محوري بالنسبة لمحور الفواصل.
                </div>
            </div>
            <div style="flex: 1; background: #ffffff; border: 2px solid #e2e8f0; border-radius: 20px; padding: 30px; box-shadow: 0 10px 25px rgba(0,0,0,0.05);">
                <div style="font-size: 24px; color: #0f172a; font-weight: 700; margin-bottom: 15px;">Le Module · المعيار (المسافة)</div>
                <div style="font-size: 38px; color: #16a34a; font-weight: 800; font-family: monospace; background: #f0fdf4; padding: 15px; border-radius: 12px;">
                    |z| = √(a² + b²)
                </div>
                <div style="font-size: 20px; color: #475569; margin-top: 15px; direction: rtl; font-family: 'Cairo', sans-serif;">
                    يمثل المسافة من نقطة الأصل O إلى النقطة M(a, b). دائماً موجب أو معدوم.
                </div>
            </div>
        </div>
        """,
        "voice_text": "توا نشوفو المرافق والمعيار. مرافق زاد هو زاد بار، ويساوي أ إلا إي ب، يعني نبدلو علامة الجزء التخيلي. والمعيار هو القيمة المطلقة لزاد، وتساوي الجذر التربيعي لـ أ قوّة اثنين مع ب قوّة اثنين. ويمثل هندسياً المسافة بين نقطة الأصل والنقطة إم في المستوى المركب."
    },
    {
        "id": "slide4",
        "title": "3. Exemple d'Application Bac · مثال تطبيقي",
        "subtitle": "تطبيق عددي مباشر على المعيار والمرافق",
        "badge": "EXEMPLE D'EXAMEN",
        "content_html": """
        <div style="background: #ffffff; border: 2px solid #0284c7; border-radius: 24px; padding: 35px; margin-top: 25px; box-shadow: 0 10px 30px rgba(2, 132, 199, 0.08); text-align: left;">
            <div style="font-size: 32px; color: #0284c7; font-weight: 800; font-family: monospace; margin-bottom: 25px; padding-bottom: 15px; border-bottom: 1px solid #e2e8f0;">
                Soit le nombre complexe : z = 3 + 4i
            </div>
            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 30px;">
                <div style="background: #f8fafc; padding: 25px; border-radius: 16px; border: 1px solid #e2e8f0;">
                    <div style="font-size: 22px; color: #475569; font-weight: 600; margin-bottom: 10px;">Calcul du Conjugué :</div>
                    <div style="font-size: 36px; color: #0f172a; font-weight: 800; font-family: monospace;">
                        z̄ = 3 - 4i
                    </div>
                </div>
                <div style="background: #f8fafc; padding: 25px; border-radius: 16px; border: 1px solid #e2e8f0;">
                    <div style="font-size: 22px; color: #475569; font-weight: 600; margin-bottom: 10px;">Calcul du Module :</div>
                    <div style="font-size: 30px; color: #16a34a; font-weight: 800; font-family: monospace;">
                        |z| = √(3² + 4²) = √(9 + 16) = √25 = 5
                    </div>
                </div>
            </div>
        </div>
        """,
        "voice_text": "ناخذو مثال تطبيقي ساهل يتعاود ديما في امتحانات الباك. إذا كان زاد يساوي ثلاثة مع أربعة إي. المرافق زاد بار يساوي ثلاثة إلا أربعة إي. والمعيار يساوي جذر ثلاثة قوّة اثنين مع أربعة قوّة اثنين، يعني جذر تسعة مع سطاش، يساوي جذر خمسة وعشرين، والنتيجة هي خمسة."
    },
    {
        "id": "slide5",
        "title": "4. À Vous de Jouer ! · دوركم في التطبيق",
        "subtitle": "واصل التمارين التفاعلية على المنصة",
        "badge": "EXERCICES ET SUITE DU COURS",
        "content_html": """
        <div style="margin-top: 40px; background: #ffffff; border: 2px solid #e2e8f0; border-radius: 24px; padding: 40px; box-shadow: 0 10px 25px rgba(0,0,0,0.05); display: inline-block;">
            <div style="font-size: 36px; color: #0284c7; font-weight: 800; font-family: 'Cairo', sans-serif; direction: rtl; margin-bottom: 20px;">
                الخطوة الموالية : حل التمارين التفاعلية
            </div>
            <div style="font-size: 22px; color: #475569; line-height: 1.8; direction: rtl; font-family: 'Cairo', sans-serif;">
                ✓ الشكل المثلثي والآسي (Forme trigonométrique & exponentielle)<br>
                ✓ معادلات الدرجة الثانية في مجموعة الأعداد المركبة ℂ<br>
                ✓ تمارين الباكالوريا الوطنية مع الإصلاح المفصل
            </div>
            <div style="margin-top: 30px; font-size: 26px; color: #16a34a; font-weight: 700; direction: rtl; font-family: 'Cairo', sans-serif;">
                بالتوفيق والنجاح لجميع تلامذة الباك في تونس ! 🇹🇳
            </div>
        </div>
        """,
        "voice_text": "وهكا نكونو كملنا المفاهيم الأساسية للدرس. توا تنجمو تعديو للتطبيق التفاعلي وتحلو التمارين بالخطوات. بالنجاح والتوفيق للناس الكل في الباك!"
    }
]

html_template = """<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="utf-8">
    <link href="https://fonts.googleapis.com/css2?family=Cairo:wght@600;700;800;900&family=Plus+Jakarta+Sans:wght@500;600;700;800&display=swap" rel="stylesheet">
    <style>
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body {
            width: 1920px;
            height: 1080px;
            background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
            font-family: 'Plus Jakarta Sans', sans-serif;
            color: #0f172a;
            display: flex;
            flex-direction: column;
            justify-content: space-between;
            padding: 80px 100px;
            overflow: hidden;
        }
        .header {
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .badge {
            background: #0284c7;
            color: white;
            padding: 8px 20px;
            border-radius: 9999px;
            font-size: 18px;
            font-weight: 700;
            letter-spacing: 0.05em;
        }
        .brand {
            font-size: 28px;
            font-weight: 800;
            color: #0284c7;
            display: flex;
            align-items: center;
            gap: 12px;
        }
        .brand span {
            color: #dc2626;
        }
        .main {
            text-align: center;
            margin: auto 0;
        }
        .title {
            font-size: 54px;
            font-weight: 800;
            color: #0f172a;
            line-height: 1.2;
            margin-bottom: 16px;
        }
        .subtitle {
            font-size: 28px;
            color: #64748b;
            font-weight: 600;
            margin-bottom: 20px;
        }
        .footer {
            display: flex;
            justify-content: space-between;
            align-items: center;
            border-top: 1px solid #e2e8f0;
            padding-top: 25px;
            color: #94a3b8;
            font-size: 20px;
            font-weight: 600;
        }
    </style>
</head>
<body>
    <div class="header">
        <div class="brand">LEARNOO<span>.TN</span> · درس</div>
        <div class="badge">__BADGE__</div>
    </div>
    <div class="main">
        <h1 class="title">__TITLE__</h1>
        <p class="subtitle">__SUBTITLE__</p>
        __CONTENT__
    </div>
    <div class="footer">
        <div>Programme Officiel · Ministère de l'Éducation Tunisienne</div>
        <div>Baccalauréat Mathématiques & Sciences</div>
    </div>
</body>
</html>
"""

print("1. Generating slides HTML and voiceover audio...")
segment_files = []

for i, slide in enumerate(slides_data):
    slide_id = slide["id"]
    html_content = html_template.replace("__BADGE__", slide["badge"])
    html_content = html_content.replace("__TITLE__", slide["title"])
    html_content = html_content.replace("__SUBTITLE__", slide["subtitle"])
    html_content = html_content.replace("__CONTENT__", slide["content_html"])
    
    html_path = os.path.join(SCRIPT_DIR, f"{slide_id}.html")
    img_path = os.path.join(SCRIPT_DIR, f"{slide_id}.png")
    audio_path = os.path.join(SCRIPT_DIR, f"{slide_id}.mp3")
    seg_video = os.path.join(SCRIPT_DIR, f"{slide_id}.mp4")

    with open(html_path, "w", encoding="utf-8") as f:
        f.write(html_content)

    # 1. Take Screenshot of Slide
    subprocess.run([
        "google-chrome", "--headless=new", "--disable-gpu",
        "--window-size=1920,1080", f"--screenshot={img_path}",
        html_path
    ], check=True)

    # 2. Generate Audio via Edge-TTS Tunisian Voice
    subprocess.run([
        os.path.expanduser("~/.local/bin/edge-tts"),
        "--voice", "ar-TN-HediNeural",
        "--text", slide["voice_text"],
        "--write-media", audio_path
    ], check=True)

    # 3. Create MP4 video segment combining still image and audio
    subprocess.run([
        "ffmpeg", "-y", "-loop", "1", "-i", img_path,
        "-i", audio_path,
        "-c:v", "libx264", "-tune", "stillimage",
        "-c:a", "aac", "-b:a", "192k", "-pix_fmt", "yuv420p",
        "-shortest", seg_video
    ], check=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)

    segment_files.append(seg_video)
    print(f"✓ Segment {i+1}/{len(slides_data)} ready: {seg_video}")

# 4. Concatenate all segments into the final master video
concat_txt = os.path.join(SCRIPT_DIR, "concat.txt")
with open(concat_txt, "w") as f:
    for seg in segment_files:
        f.write(f"file '{seg}'\n")

final_output = os.path.join(SCRIPT_DIR, "lesson_nombres_complexes_tn.mp4")
print("\n2. Concatenating master video...")
subprocess.run([
    "ffmpeg", "-y", "-f", "concat", "-safe", "0",
    "-i", concat_txt,
    "-c", "copy", final_output
], check=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)

size_mb = os.path.getsize(final_output) / (1024 * 1024)
print(f"\n🎉 MASTER VIDEO COMPLETE: {final_output} ({size_mb:.2f} MB)")
