"""Unreviewed technical slide content shared by both narration variants."""

SLIDES = [{'id': 'slide1',
  'title': 'Nombres Complexes · الأعداد المركبة',
  'subtitle': 'Forme algébrique, conjugué et module · اختبار تقني غير معتمد',
  'badge': 'PROTOTYPE NON VALIDÉ · غير معتمد',
  'content_html': '\n'
                  '        <div style="margin-top: 40px; background: rgba(2, 132, 199, 0.08); border: 2px '
                  'solid rgba(2, 132, 199, 0.2); border-radius: 24px; padding: 36px 48px; display: '
                  'inline-block;">\n'
                  '            <div style="font-size: 38px; color: #0284c7; font-weight: 800; font-family: '
                  '\'Cairo\', sans-serif; direction: rtl;">\n'
                  '                شرح مبسّط بالدارجة التونسية للمفاهيم الأساسية\n'
                  '            </div>\n'
                  '            <div style="font-size: 24px; color: #475569; margin-top: 14px; font-family: '
                  '\'Plus Jakarta Sans\', sans-serif;">\n'
                  '                Section Mathématiques & Sciences Expérimentales\n'
                  '            </div>\n'
                  '        </div>\n'
                  '        ',
  'voice_text': 'عسلامة ومرحباً بيكم تلامذة الباكالوريا في تونس. في هذا الفيديو القصير باش نفهمو مع بعضنا '
                'درس الأعداد المركبة، الشكل الجبري والمعيار، بالدارجة التونسية وبأبسط طريقة.',
  'msa_text': 'مرحباً بكم. هذا اختبار تقني غير معتمد حول les nombres complexes. نعرض الشكل الجبري والمرافق '
              'والمعيار. يحتاج المحتوى إلى مراجعة من أستاذ قبل استعماله في التعليم.'},
 {'id': 'slide2',
  'title': '1. La Forme Algébrique · الشكل الجبري',
  'subtitle': 'تعريف العدد المركب والقاعدة الذهبية',
  'badge': 'DÉFINITION FONDAMENTALE',
  'content_html': '\n'
                  '        <div style="display: flex; gap: 40px; justify-content: center; margin-top: '
                  '30px;">\n'
                  '            <div style="flex: 1; background: #ffffff; border: 2px solid #e2e8f0; '
                  'border-radius: 20px; padding: 30px; box-shadow: 0 10px 25px rgba(0,0,0,0.05);">\n'
                  '                <div style="font-size: 24px; color: #0f172a; font-weight: 700; '
                  'margin-bottom: 15px;">La Règle d\'or · القاعدة الأساسية</div>\n'
                  '                <div style="font-size: 42px; color: #dc2626; font-weight: 800; '
                  'font-family: monospace; background: #fef2f2; padding: 15px; border-radius: 12px;">\n'
                  '                    i² = -1\n'
                  '                </div>\n'
                  '                <div style="font-size: 20px; color: #64748b; margin-top: 15px; direction: '
                  'rtl; font-family: \'Cairo\', sans-serif;">\n'
                  '                    عدد جديد يسمح بحل المعادلات التي ليس لها حل في ℝ.\n'
                  '                </div>\n'
                  '            </div>\n'
                  '            <div style="flex: 1; background: #ffffff; border: 2px solid #e2e8f0; '
                  'border-radius: 20px; padding: 30px; box-shadow: 0 10px 25px rgba(0,0,0,0.05);">\n'
                  '                <div style="font-size: 24px; color: #0f172a; font-weight: 700; '
                  'margin-bottom: 15px;">Forme Algébrique · الشكل الجبري</div>\n'
                  '                <div style="font-size: 42px; color: #0284c7; font-weight: 800; '
                  'font-family: monospace; background: #f0f9ff; padding: 15px; border-radius: 12px;">\n'
                  '                    z = a + i b\n'
                  '                </div>\n'
                  '                <div style="font-size: 20px; color: #334155; margin-top: 15px; direction: '
                  'rtl; font-family: \'Cairo\', sans-serif; line-height: 1.6;">\n'
                  '                    • <b>a</b> هو الجزء الحقيقي (Partie réelle): Re(z) = a<br>\n'
                  '                    • <b>b</b> هو الجزء التخيلي (Partie imaginaire): Im(z) = b\n'
                  '                </div>\n'
                  '            </div>\n'
                  '        </div>\n'
                  '        ',
  'voice_text': 'نبداو بالقاعدة الأساسية. في مجموعة الأعداد المركبة، فما عدد جديد اسمو إي، بحيث إي قوّة '
                'اثنين تساوي إلا واحد. وأي عدد مركب زاد يتكتب على الشكل الجبري: زاد يساوي أ مع إي ب. أ هو '
                'الجزء الحقيقي، وب هو الجزء التخيلي.',
  'msa_text': 'لكل عدد مركب z شكل جبري: z égale a plus i b، حيث a و b عددان حقيقيان و i carré égale moins '
              'un. a هو la partie réelle و b هو la partie imaginaire.'},
 {'id': 'slide3',
  'title': '2. Conjugué et Module · المرافق والمعيار',
  'subtitle': 'المفاهيم الهندسية والمسافة في المستوى المركب',
  'badge': 'GÉOMÉTRIE DANS LE PLAN COMPLEXE',
  'content_html': '\n'
                  '        <div style="display: flex; gap: 40px; justify-content: center; margin-top: '
                  '30px;">\n'
                  '            <div style="flex: 1; background: #ffffff; border: 2px solid #e2e8f0; '
                  'border-radius: 20px; padding: 30px; box-shadow: 0 10px 25px rgba(0,0,0,0.05);">\n'
                  '                <div style="font-size: 24px; color: #0f172a; font-weight: 700; '
                  'margin-bottom: 15px;">Le Conjugué · مرافق عدد مركب</div>\n'
                  '                <div style="font-size: 38px; color: #0284c7; font-weight: 800; '
                  'font-family: monospace; background: #f0f9ff; padding: 15px; border-radius: 12px;">\n'
                  '                    z̄ = a - i b\n'
                  '                </div>\n'
                  '                <div style="font-size: 20px; color: #475569; margin-top: 15px; direction: '
                  'rtl; font-family: \'Cairo\', sans-serif;">\n'
                  '                    نبدلو علامة الجزء التخيلي فقط. هندسياً هو تناظر محوري بالنسبة لمحور '
                  'الفواصل.\n'
                  '                </div>\n'
                  '            </div>\n'
                  '            <div style="flex: 1; background: #ffffff; border: 2px solid #e2e8f0; '
                  'border-radius: 20px; padding: 30px; box-shadow: 0 10px 25px rgba(0,0,0,0.05);">\n'
                  '                <div style="font-size: 24px; color: #0f172a; font-weight: 700; '
                  'margin-bottom: 15px;">Le Module · المعيار (المسافة)</div>\n'
                  '                <div style="font-size: 38px; color: #16a34a; font-weight: 800; '
                  'font-family: monospace; background: #f0fdf4; padding: 15px; border-radius: 12px;">\n'
                  '                    |z| = √(a² + b²)\n'
                  '                </div>\n'
                  '                <div style="font-size: 20px; color: #475569; margin-top: 15px; direction: '
                  'rtl; font-family: \'Cairo\', sans-serif;">\n'
                  '                    يمثل المسافة من نقطة الأصل O إلى النقطة M(a, b). دائماً موجب أو '
                  'معدوم.\n'
                  '                </div>\n'
                  '            </div>\n'
                  '        </div>\n'
                  '        ',
  'voice_text': 'توا نشوفو المرافق والمعيار. مرافق زاد هو زاد بار، ويساوي أ إلا إي ب، يعني نبدلو علامة الجزء '
                'التخيلي. والمعيار هو القيمة المطلقة لزاد، وتساوي الجذر التربيعي لـ أ قوّة اثنين مع ب قوّة '
                'اثنين. ويمثل هندسياً المسافة بين نقطة الأصل والنقطة إم في المستوى المركب.',
  'msa_text': 'المرافق le conjugué للعدد a plus i b هو a moins i b. أما المعيار le module فهو الجذر التربيعي '
              'لـ a carré plus b carré. وهو عدد حقيقي موجب أو منعدم، يمثل المسافة إلى نقطة الأصل.'},
 {'id': 'slide4',
  'title': '3. Exemple technique · مثال تقني',
  'subtitle': 'تطبيق عددي مباشر على المعيار والمرافق',
  'badge': 'EXEMPLE NON VALIDÉ',
  'content_html': '\n'
                  '        <div style="background: #ffffff; border: 2px solid #0284c7; border-radius: 24px; '
                  'padding: 35px; margin-top: 25px; box-shadow: 0 10px 30px rgba(2, 132, 199, 0.08); '
                  'text-align: left;">\n'
                  '            <div style="font-size: 32px; color: #0284c7; font-weight: 800; font-family: '
                  'monospace; margin-bottom: 25px; padding-bottom: 15px; border-bottom: 1px solid '
                  '#e2e8f0;">\n'
                  '                Soit le nombre complexe : z = 3 + 4i\n'
                  '            </div>\n'
                  '            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 30px;">\n'
                  '                <div style="background: #f8fafc; padding: 25px; border-radius: 16px; '
                  'border: 1px solid #e2e8f0;">\n'
                  '                    <div style="font-size: 22px; color: #475569; font-weight: 600; '
                  'margin-bottom: 10px;">Calcul du Conjugué :</div>\n'
                  '                    <div style="font-size: 36px; color: #0f172a; font-weight: 800; '
                  'font-family: monospace;">\n'
                  '                        z̄ = 3 - 4i\n'
                  '                    </div>\n'
                  '                </div>\n'
                  '                <div style="background: #f8fafc; padding: 25px; border-radius: 16px; '
                  'border: 1px solid #e2e8f0;">\n'
                  '                    <div style="font-size: 22px; color: #475569; font-weight: 600; '
                  'margin-bottom: 10px;">Calcul du Module :</div>\n'
                  '                    <div style="font-size: 30px; color: #16a34a; font-weight: 800; '
                  'font-family: monospace;">\n'
                  '                        |z| = √(3² + 4²) = √(9 + 16) = √25 = 5\n'
                  '                    </div>\n'
                  '                </div>\n'
                  '            </div>\n'
                  '        </div>\n'
                  '        ',
  'voice_text': 'ناخذو مثال تقني: زاد يساوي ثلاثة مع أربعة إي. المرافق يساوي ثلاثة إلا أربعة إي. والمعيار '
                'يساوي جذر تسعة مع سطاش، والنتيجة خمسة. المثال هذا موش درس مصادق عليه.',
  'msa_text': 'في هذا المثال التقني، z égale trois plus quatre i. المرافق هو trois moins quatre i. والمعيار '
              'يساوي الجذر التربيعي لتسعة زائد ستة عشر، أي خمسة. المثال غير معتمد تربوياً.'},
 {'id': 'slide5',
  'title': '4. Fin du prototype · نهاية الاختبار',
  'subtitle': 'Contenu non validé — aucune publication pédagogique',
  'badge': 'REVUE ENSEIGNANTE REQUISE',
  'content_html': '<p style="font-size:36px;line-height:1.8">Ce prototype technique attend une revue '
                  'enseignante.<br>هذا اختبار تقني غير معتمد، يحتاج إلى مراجعة من أستاذ.</p>',
  'voice_text': 'هنا يوفى الاختبار التقني. المحتوى هذا مازال موش مصادق عليه ويلزمو مراجعة من أستاذ قبل '
                'الاستعمال في التعليم.',
  'msa_text': 'هنا ينتهي الاختبار التقني. المحتوى غير معتمد ويحتاج إلى مراجعة من أستاذ قبل استعماله في '
              'التعليم.'}]

HTML_TEMPLATE = '<!DOCTYPE html>\n<html lang="fr">\n<head>\n    <meta charset="utf-8">\n    <style>\n        * { box-sizing: border-box; margin: 0; padding: 0; }\n        body {\n            width: 1920px;\n            height: 1080px;\n            background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);\n            font-family: \'Plus Jakarta Sans\', sans-serif;\n            color: #0f172a;\n            display: flex;\n            flex-direction: column;\n            justify-content: space-between;\n            padding: 80px 100px;\n            overflow: hidden;\n        }\n        .header {\n            display: flex;\n            justify-content: space-between;\n            align-items: center;\n        }\n        .badge {\n            background: #0284c7;\n            color: white;\n            padding: 8px 20px;\n            border-radius: 9999px;\n            font-size: 18px;\n            font-weight: 700;\n            letter-spacing: 0.05em;\n        }\n        .brand {\n            font-size: 28px;\n            font-weight: 800;\n            color: #0284c7;\n            display: flex;\n            align-items: center;\n            gap: 12px;\n        }\n        .brand span {\n            color: #dc2626;\n        }\n        .main {\n            text-align: center;\n            margin: auto 0;\n        }\n        .title {\n            font-size: 54px;\n            font-weight: 800;\n            color: #0f172a;\n            line-height: 1.2;\n            margin-bottom: 16px;\n        }\n        .subtitle {\n            font-size: 28px;\n            color: #64748b;\n            font-weight: 600;\n            margin-bottom: 20px;\n        }\n        .footer {\n            display: flex;\n            justify-content: space-between;\n            align-items: center;\n            border-top: 1px solid #e2e8f0;\n            padding-top: 25px;\n            color: #94a3b8;\n            font-size: 20px;\n            font-weight: 600;\n        }\n    </style>\n</head>\n<body>\n    <div class="header">\n        <div class="brand">LEARNOO<span>.TN</span> · درس</div>\n        <div class="badge">__BADGE__</div>\n    </div>\n    <div class="main">\n        <h1 class="title">__TITLE__</h1>\n        <p class="subtitle">__SUBTITLE__</p>\n        __CONTENT__\n    </div>\n    <div class="footer">\n        <div>Prototype non validé · Revue enseignante requise</div>\n        <div>Test local — non publié</div>\n    </div>\n</body>\n</html>\n'
