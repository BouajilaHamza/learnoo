import os
import asyncio
import edge_tts

SCRIPTS = [
    {
        "slide": 1,
        "title": "Introduction aux Nombres Complexes",
        "text": "مرحباً بكم. نتناول اليوم درساً محورياً في برنامج البكالوريا: les nombres complexes. نشأ هذا المفهوم الرياضي عندما واجه العلماء معادلات لا تقبل أي حل في مجموعة الأعداد الحقيقية R، مثل المعادلة x carré plus 1 égale 0. من هنا تم إدخال العدد التخيلي i، بحيث i carré égale moins 1."
    },
    {
        "slide": 2,
        "title": "La Forme Algébrique",
        "text": "لكل عدد مركب z تمثيل جبري يسمى la forme algébrique: z égale a plus i b. يسمى a الجزء الحقيقي، la partie réelle de z. بينما يسمى b الجزء التخيلي، la partie imaginaire de z. تذكروا دائماً أن b هو عدد حقيقي يمثل معامل i."
    },
    {
        "slide": 3,
        "title": "Conjugué et Module",
        "text": "من المفاهيم الأساسية: المرافق le conjugué، ونرمز له بـ z barre، حيث نغير إشارة الجزء التخيلي: a moins i b. أما المعيار le module de z، فهو المسافة في المستوي المركب، ويحسب بالجذر التربيعي لـ a carré plus b carré. المعيار دائماً عدد حقيقي موجب."
    },
    {
        "slide": 4,
        "title": "Forme Trigonométrique et Exponentielle",
        "text": "في الهندسة المستوية، يمكن التعبير عن العدد المركب بصيغتين: la forme trigonométrique و la forme exponentielle. إذا كان r هو le module، و thêta هو l'argument، فإن z يكتب r facteur de cosinus thêta plus i sinus thêta، أو بصيغة أويلر: r exponentielle i thêta. هذه الصيغة تسهل كثيراً حساب الجداء وقوى الأعداد المركبة."
    },
    {
        "slide": 5,
        "title": "Application Type Bac",
        "text": "نختم بتطبيق نموذجي من مواضيع البكالوريا: ليكن z égale 1 plus i racine de 3. بحساب le module نجد r égale 2، و l'argument يوافق pi sur 3. بالتالي la forme exponentielle هي: 2 exponentielle i pi sur 3. بهذا المنهج المنظم تضمنون النجاح في الامتحان."
    }
]

async def generate():
    os.makedirs("/tmp/learnoo_video/msa", exist_ok=True)
    voice = "ar-TN-HediNeural"
    for item in SCRIPTS:
        s_num = item["slide"]
        out_mp3 = f"/tmp/learnoo_video/msa/slide{s_num}_msa.mp3"
        print(f"Generating audio for slide {s_num}...")
        communicate = edge_tts.Communicate(item["text"], voice, rate="+0%")
        await communicate.save(out_mp3)
        print(f"✓ Saved {out_mp3}")

if __name__ == "__main__":
    asyncio.run(generate())
