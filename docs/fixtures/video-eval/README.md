# Évaluation de la Génération Vidéo et Audio Automatisée pour Learnoo

Ce document consigne les tests empiriques et l'évaluation technique de la génération automatisée de vidéos de cours pour les programmes du Baccalauréat tunisien (matière : Mathématiques, chapitre : Nombres complexes).

Deux prototypes vidéo complets (1080p Full HD) et leurs pistes audio respectives ont été générés et sont fournis dans `assets/fixtures/video/`.

---

## 1. Prototypes Fournis

| Version | Piste Audio | Fichier Vidéo | Durée | Approche Linguistique | Évaluation Pédagogique |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **V1: Derja Tunisienne** | `cours_nombres_complexes_derja.mp3` | `lesson_nombres_complexes_derja.mp4` | 91s | Derja familière (`ar-TN-HediNeural`) | **Insuffisante (Rejetée)**: Trop informelle, intonation artificielle, inadéquate pour un cours académique formel. |
| **V2: Arabe Standard (MSA) + Français** | `cours_nombres_complexes_msa.mp3` | `cours_nombres_complexes_msa_video.mp4` | 111s | Arabe académique moderne (الفصحى) + termes techniques en français | **Supérieure**: Posture pédagogique conforme aux lycées pilotes tunisiens, prononciation claire. |

---

## 2. Pipeline Technique Utilisé

Les scripts de génération automatisée sont disponibles dans `docs/fixtures/video-eval/scripts/` :
1. **`generate_lesson_video.py`** : Génération des diapositives HTML avec rendu mathématique KaTeX, capture d'écran 1920x1080 sans tête (headless Chromium), synthèse vocale neural TTS (`edge-tts`) et assemblage segment par segment.
2. **`generate_msa_audio.py`** : Synthèse de la piste en Arabe Standard Moderne avec intégration des termes mathématiques français latins (*la forme algébrique*, *la partie réelle*, *la partie imaginaire*, *le conjugué*, *le module*, *l'argument*, *la forme exponentielle*).
3. **`render_msa_video.py`** : Concaténation de la piste audio maîtresse et multiplexage vidéo H.264 / AAC 1080p via `ffmpeg`.

---

## 3. Limites de la Génération Vidéo Automatisée par IA

Bien que le pipeline produise un rendu Full HD synchronisé en quelques minutes, l'évaluation qualitative montre des faiblesses structurelles pour une plateforme éducative rigoureuse :

1. **Absence de tracé dynamique (Tableau / Manim)** :
   - Les mathématiques de niveau Baccalauréat s'expliquent par le geste et l'écriture pas-à-pas (décomposition d'une équation, calcul progressif du discriminant ou du module).
   - Les diapositives statiques créent une charge cognitive excessive pour l'élève.

2. **Problématique du bilinguisme franco-arabe (Code-switching)** :
   - La voix neuronale arabe prononçant des termes latins français produit parfois des heuristiques phonétiques imparfaites sans transcription phonétique rigoureuse.
   - Les dialectes régionaux (Derja tunisienne brute) manquent de modèles neuronaux académiques calibrés pour l'enseignement secondaire.

3. **Rythme et pauses cognitives** :
   - Un enseignant ajuste son débit selon la difficulté du concept et marque des silences intentionnels. Le TTS automatisé applique un tempo uniforme qui nuit à la rétention.

---

## 4. Recommandations pour Learnoo

1. **Conserver la vidéo humaine sur YouTube** :
   - Conforme au principe défini dans `docs/ARCHITECTURE.md` ("*YouTube delivers video*").
   - Privilégier les enregistrements d'enseignants tunisiens réels et qualifiés.

2. **Usage de la synthèse vocale / vidéo IA** :
   - Utiliser l'audio automatisé uniquement pour l'accessibilité (lecture audio pour malvoyants du contenu texte/KaTeX des fiches de cours).
   - Ne pas substituer les vidéos pédagogiques interactives par de simples montages automatisés de diapositives statiques.
