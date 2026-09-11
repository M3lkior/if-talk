# Alignement du rendu des slides sur l'identité zatsit.fr — plan d'implémentation

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Faire lire le deck `impact-framework` comme la nouvelle version de zatsit.fr, en clair et en sombre, sans toucher au châssis ni aux autres talks.

**Architecture :** La palette runtime vit dans le moteur (`styles/tokens.css` pour le clair, le bloc `.dark` de `styles/demoit.src.css` pour le sombre) ; le `tokens.css` du talk en reste le reflet, parce qu'importé en `theme(reference)` il n'émet rien. Le décor est une propriété de la scène — deux pseudo-éléments sur `.stage` pilotés par deux jetons d'opacité — donc aucun markup de slide ne bouge. La cover inline le SVG du logo pour que le dégradé du thème puisse le remplir, et tire son titre de `talk.yml`.

**Tech Stack :** Go 1.27, Tailwind CSS v4 (CLI standalone épinglé 4.3.3), templates `html/template`, goldmark.

**Spec :** `docs/superpowers/specs/2026-09-11-zatsit-design-system-design.md`

## Global Constraints

- **Toute commande Go exige `GOCACHE="$TMPDIR/go-build"`.** Le bac à sable refuse d'écrire dans `~/Library/Caches/go-build` et la commande échoue sur `operation not permitted`.
- **Les deux feuilles de style sont des artefacts commités qu'aucun `go build` ne régénère.** Une classe Tailwind écrite dans une slide n'existe pas tant que `hack/css.sh impact-framework` n'a pas tourné ; une classe écrite dans une layout ou dans `index.tmpl.html` n'existe pas tant que `hack/css.sh engine` n'a pas tourné. Tout commit qui touche à `styles/` ou à un fichier scanné inclut le CSS régénéré.
- **`styles/tokens.css` et `impact-framework/.demoit/tokens.css` ne doivent contenir que des blocs `@theme`.** Tailwind rejette un fichier importé en `theme(reference)` qui porte une autre règle.
- **Les fichiers `_test.go` ne portent pas l'en-tête Apache.** C'est la convention établie ici ; aucun fichier de test du dépôt n'en a, et `goheader` n'a pas de template configuré dans `golangci.yml`.
- **Les messages de commit suivent le conventional commit `type(scope): message`.** Les scopes utilisés dans l'historique récent : `impact-framework`, `specs`, et le nom du paquet pour le moteur.
- **Un talk doit rendre sans réseau.** Le réseau ne sert qu'au build (CLI Tailwind, récupération des fonts) et jamais pendant une présentation.
- **`sample/` ne doit rien voir de ce chantier.** Il ne déclare pas de `talk.yml`, rend en `.no-stage` et n'a pas de thème sombre. C'est un critère de vérification, pas une supposition.

---

### Task 1 : La palette

**Files:**
- Modify: `styles/tokens.css`
- Modify: `styles/demoit.src.css:123-173` (le bloc `:root` des jetons de lisibilité et le bloc `.dark`)
- Modify: `impact-framework/.demoit/tokens.css`
- Regenerate: `handlers/resources/demoit.css`, `impact-framework/.demoit/tailwind.css`

**Interfaces:**
- Produces: les jetons `--color-accent-1`, `--color-accent-2`, `--color-accent-3`, `--color-surface-raised`, `--logo-gradient-start`, `--logo-gradient-mid`, `--logo-gradient-end`, `--color-text-gradient`, `--decor-dots`, `--decor-mesh`. Les tâches 3, 4 et 7 les lisent.

**Écart assumé avec la spec :** son tableau range `--color-text-gradient` parmi les jetons `@theme`. Il n'y va pas. `@theme` interprète le préfixe `--color-*` comme le namespace des couleurs et en génère des utilitaires (`bg-text-gradient`, `text-text-gradient`) qui seraient cassés, la valeur étant un `linear-gradient` et non une couleur. Il est donc déclaré comme custom property ordinaire dans le `:root` de `demoit.src.css`, à côté de `--rule` et des `--dm-*`, et permuté dans `.dark`. Même endroit, même effet, sans utilitaire bancal.

- [ ] **Step 1 : Réécrire la palette claire du moteur**

Dans `styles/tokens.css`, remplacer le bloc des couleurs (les lignes `--color-main` à `--color-outline`) et `--radius-card` par :

```css
    /* Palette de la maison, alignée sur corporate/src/styles/global.css du
       dépôt zats-websites. Un talk la reflète dans son propre tokens.css pour
       que ses utilitaires se génèrent contre le bon vocabulaire, mais c'est
       bien ce fichier-ci qui l'émet : un tokens.css de talk est importé en
       theme(reference) et ne produit aucune variable. */
    --color-main: #0f15fd;
    --color-main-dark: #f1be51;

    --color-fg: #1c1e21;
    --color-fg-muted: #475569;
    --color-surface: #ffffff;
    --color-surface-raised: #f8fafc;
    --color-outline: #cbd5e1;

    --color-accent-1: #06b6d4;
    --color-accent-2: #3b82f6;
    --color-accent-3: #6366f1;

    --logo-gradient-start: #0F15FD;
    --logo-gradient-mid: #0a0ecc;
    --logo-gradient-end: #020466;

    /* Le létterboxage autour de la scène. Pas d'équivalent sur le site, et
       c'est voulu qu'il reste neutre. */
    --color-void: #000000;
```

et plus bas dans le même bloc `@theme` :

```css
    --radius-card: 0.75rem;
```

Le commentaire existant qui dit « A talk overrides --color-main in its own .demoit/tokens.css » est faux et doit disparaître avec ces lignes : c'est précisément ce que la mesure a infirmé.

- [ ] **Step 2 : Refléter la palette claire dans le talk**

Dans `impact-framework/.demoit/tokens.css`, poser **exactement les mêmes valeurs** que le Step 1 pour `--color-main`, `--color-main-dark`, `--color-fg`, `--color-fg-muted`, `--color-surface`, `--color-surface-raised`, `--color-outline`, `--color-accent-1/2/3`, `--logo-gradient-start/mid/end`, `--color-void` et `--radius-card`.

Remplacer aussi le commentaire de tête du fichier, qui affirme que la palette du talk vit ici, par :

```css
/*
Reflet de styles/tokens.css : le vocabulaire contre lequel les utilitaires de
ce talk sont générés.

Ce fichier n'émet rien. tailwind.src.css l'importe en `theme(reference)`, ce
qui donne les noms au générateur sans produire une seule variable -- vérifiable
dans tailwind.css, qui ne contient aucun bloc :root. Les valeurs qui
s'appliquent réellement sont celles du :root et du .dark de
handlers/resources/demoit.css.

Faire diverger ce fichier de celui du moteur ne repeint donc rien : ça ne fait
que générer des utilitaires dont le repli est mort. Les deux se modifient
ensemble.

Doit contenir des blocs @theme et rien d'autre.
*/
```

- [ ] **Step 3 : Ajouter le dégradé de texte et les jetons de décor au `:root` du moteur**

Dans `styles/demoit.src.css`, à la fin du bloc `:root` existant (celui qui commence par `--rule: 1px;`), avant l'accolade fermante :

```css
    /* Un linear-gradient, pas une couleur : déclaré ici plutôt que dans
       @theme, où le préfixe --color-* génèrerait des utilitaires de couleur
       cassés. Lu par le traitement `em` de slide-hero. */
    --color-text-gradient: linear-gradient(to bottom, #0F15FD 0%, #020466 100%);

    /* Opacité des deux couches de décor de la scène, une par couche : le site
       baisse son mesh en sombre et laisse sa grille de points inchangée. */
    --decor-dots: 0.08;
    --decor-mesh: 0.30;
```

- [ ] **Step 4 : Basculer le bloc `.dark` sur l'or**

Dans `styles/demoit.src.css`, remplacer le contenu du bloc `.dark` par :

```css
.dark {
    --color-fg: #e3e3e3;
    --color-fg-muted: #94a3b8;
    --color-surface: #1b1b1d;
    --color-surface-raised: rgba(255, 255, 255, 0.05);
    --color-outline: #334155;
    --color-main: var(--color-main-dark);

    /* Le site inverse la teinte en sombre : primaire or, accents ambre,
       rouge, orange. C'est la signature du nouveau zatsit, pas un
       éclaircissement du bleu. */
    --color-accent-1: #f59e0b;
    --color-accent-2: #ef4444;
    --color-accent-3: #f97316;

    --logo-gradient-start: #f1be51;
    --logo-gradient-mid: #f59e0b;
    --logo-gradient-end: #ef4444;

    --color-text-gradient: linear-gradient(135deg, #f1be51, #f59e0b, #ef4444);

    /* Le mesh baisse en sombre, comme sur le site. */
    --decor-mesh: 0.16;

    --dm-window-bg: #1c1d21;
    --dm-chrome-bg: linear-gradient(to bottom, #2b2d33 0%, #23252a 100%);
    --dm-chrome-border: #334155;
    --dm-tabs-bg: #23252a;
    --dm-tab-bg: #2b2d33;
    --dm-tab-fg: #e3e3e3;
    --dm-code-fg: #e3e3e3;
    --dm-code-selection: #2f4a7a;
}
```

`--color-main` garde son indirection par `--color-main-dark`, qui existe déjà. Les autres jetons sont réassignés directement : leur inventer un compagnon `-dark` coûterait six niveaux d'indirection pour un gain nul, le moteur portant de toute façon les deux palettes.

- [ ] **Step 5 : Régénérer les deux feuilles**

```bash
bash hack/css.sh engine
bash hack/css.sh impact-framework
```

- [ ] **Step 6 : Vérifier que la palette est bien émise**

```bash
grep -c "0f15fd" handlers/resources/demoit.css
grep -o -- "--color-main-dark: *[^;]*" handlers/resources/demoit.css
grep -o -- "--radius-card: *[^;]*" handlers/resources/demoit.css
grep -o -- "--decor-mesh: *[^;]*" handlers/resources/demoit.css
```

Attendu : au moins une occurrence de `0f15fd` ; `--color-main-dark: #f1be51` ; `--radius-card: 0.75rem` ; **deux** valeurs de `--decor-mesh` (`0.30` dans `:root`, `0.16` dans `.dark`).

- [ ] **Step 7 : Vérifier que rien de Go n'a bougé**

```bash
GOCACHE="$TMPDIR/go-build" go test ./...
```

Attendu : `ok` pour `deck`, `deck/directive`, `handlers`, `highlight`, `vscode`.

- [ ] **Step 8 : Commit**

```bash
git add styles/ handlers/resources/demoit.css impact-framework/.demoit/tokens.css impact-framework/.demoit/tailwind.css
git commit -m "feat(impact-framework): aligner la palette sur la nouvelle identité zatsit"
```

---

### Task 2 : Poppins 600

**Files:**
- Create: `impact-framework/.demoit/fonts/poppins-600-latin.woff2`
- Create: `impact-framework/.demoit/fonts/poppins-600-latin-ext.woff2`
- Modify: `impact-framework/.demoit/style.css` (bloc des `@font-face`)

**Interfaces:**
- Produces: le poids 600 de Poppins, que `font-semibold` utilise. Les tâches 4 et 7 s'en servent.

Le site charge 400 / 600 / 700 ; le talk n'embarque que 100 / 400 / 500 / 700. Sans le 600, `font-semibold` retombe silencieusement sur une synthèse du navigateur ou sur le 500.

- [ ] **Step 1 : Récupérer les deux fichiers**

Les URLs sont celles de Google Fonts pour Poppins v24, poids 600, vérifiées contre leurs `unicode-range` — le fichier `…Z1xlFd2JQEk` porte le sous-ensemble latin, le fichier `…Z1JlFd2JQEl8qw` le latin-ext, exactement les mêmes plages que les huit fichiers déjà présents.

```bash
curl -sSfL --retry 3 -o impact-framework/.demoit/fonts/poppins-600-latin.woff2 \
  "https://fonts.gstatic.com/s/poppins/v24/pxiByp8kv8JHgFVrLEj6Z1xlFd2JQEk.woff2"
curl -sSfL --retry 3 -o impact-framework/.demoit/fonts/poppins-600-latin-ext.woff2 \
  "https://fonts.gstatic.com/s/poppins/v24/pxiByp8kv8JHgFVrLEj6Z1JlFd2JQEl8qw.woff2"
```

- [ ] **Step 2 : Vérifier que ce sont bien des woff2 non vides**

```bash
file impact-framework/.demoit/fonts/poppins-600-*.woff2
ls -l impact-framework/.demoit/fonts/poppins-600-*.woff2
```

Attendu : `Web Open Font Format (Version 2)` pour les deux, et des tailles du même ordre que les autres poids (quelques dizaines de Ko), jamais 0.

**Si la récupération échoue** — réseau filtré, URL périmée : le repli est de ne pas créer ces fichiers et de mapper `font-semibold` sur le 500 déjà présent, en ajoutant dans `style.css` un `@font-face` de poids 600 pointant sur les fichiers 500. Noter le repli dans le message de commit ; c'est visible mais pas bloquant.

- [ ] **Step 3 : Déclarer les deux `@font-face`**

Dans `impact-framework/.demoit/style.css`, après le couple de poids 500 et avant celui du 700, en respectant l'ordre croissant déjà en place :

```css
@font-face {
    font-family: "Poppins";
    font-style: normal;
    font-weight: 600;
    font-display: swap;
    src: url("/fonts/poppins-600-latin.woff2") format("woff2");
    unicode-range: U+0000-00FF, U+0131, U+0152-0153, U+02BB-02BC, U+02C6, U+02DA, U+02DC, U+0304, U+0308, U+0329, U+2000-206F, U+20AC, U+2122, U+2191, U+2193, U+2212, U+2215, U+FEFF, U+FFFD;
}

@font-face {
    font-family: "Poppins";
    font-style: normal;
    font-weight: 600;
    font-display: swap;
    src: url("/fonts/poppins-600-latin-ext.woff2") format("woff2");
    unicode-range: U+0100-02BA, U+02BD-02C5, U+02C7-02CC, U+02CE-02D7, U+02DD-02FF, U+0304, U+0308, U+0329, U+1D00-1DBF, U+1E00-1E9F, U+1EF2-1EFF, U+2020, U+20A0-20AB, U+20AD-20C0, U+2113, U+2C60-2C7F, U+A720-A7FF;
}
```

Les deux plages sont copiées des `@font-face` voisins du même fichier : le sous-ensemble latin n'a pas de latin étendu et le latin-ext n'a pas les lettres de base, donc le français a besoin des deux.

- [ ] **Step 4 : Commit**

```bash
git add impact-framework/.demoit/fonts/ impact-framework/.demoit/style.css
git commit -m "feat(impact-framework): embarquer Poppins 600, le poids semi-gras du site"
```

---

### Task 3 : Le décor de la scène

**Files:**
- Modify: `styles/demoit.src.css` (après le bloc `.stage`, et le bloc `:root[data-display="projector"]`)
- Regenerate: `handlers/resources/demoit.css`

**Interfaces:**
- Consumes: `--decor-dots`, `--decor-mesh`, `--color-main`, `--color-accent-1/2/3` (tâche 1).
- Produces: les pseudo-éléments `.stage::before` et `.stage::after`. La tâche 4 en règle l'intensité.

- [ ] **Step 1 : Écrire le décor**

Dans `styles/demoit.src.css`, juste après le bloc `.stage { … }` :

```css
/*
Le décor, repris de zatsit.fr : la grille de points de .dot-grid et le mesh
gradient de .hero-gradient.

Sur la scène plutôt que sur les slides, et c'est le point : aucun markup de
slide ne le mentionne, aucune layout n'a à le connaître, et il ne peut décaler
aucune géométrie. .stage est déjà overflow:hidden, donc rien ne déborde.

Les blobs animés du site sont délibérément absents : ce deck fait tourner des
tty et des iframes live pendant trois quarts d'heure, et un blur(80px) animé en
permanence à côté est un coût permanent pour un effet que personne ne regarde.
*/
.stage::before,
.stage::after {
    content: "";
    position: absolute;
    inset: 0;
    /*
    Négatif, et c'est ce qui évite de toucher au reste. Un enfant en z-index
    négatif peint au-dessus du fond de son contexte d'empilement et en dessous
    du contenu en flux : le décor se glisse donc entre le fond de la scène et
    les slides sans qu'aucun enfant n'ait à être repositionné.

    À z-index: 0 il aurait fallu remonter les enfants avec un
    `.stage > * { position: relative }` -- qui écraserait le
    `position: absolute` de .slide-footer, une règle hors couche l'emportant
    sur un @utility en @layer utilities. Le pied de slide serait remonté dans
    le flux, la barre de progression avec.
    */
    z-index: -1;
    pointer-events: none;
}

.stage::before {
    background-image: radial-gradient(var(--color-main) 1px, transparent 1px);
    background-size: 32px 32px;
    opacity: var(--decor-dots);
}

.stage::after {
    background:
        radial-gradient(ellipse 80% 50% at 20% -20%, var(--color-accent-1), transparent 50%),
        radial-gradient(ellipse 60% 40% at 80% 0%, var(--color-main), transparent 50%),
        radial-gradient(ellipse 50% 60% at 0% 100%, var(--color-accent-3), transparent 50%),
        radial-gradient(ellipse 70% 50% at 100% 80%, var(--color-accent-2), transparent 50%);
    opacity: var(--decor-mesh);
    mask-image: linear-gradient(to bottom, black 60%, transparent 100%);
    -webkit-mask-image: linear-gradient(to bottom, black 60%, transparent 100%);
}

```

Et, dans le bloc `.stage { … }` lui-même, ajouter une ligne :

```css
    /* La scène doit créer un contexte d'empilement en toutes circonstances,
       pour que le z-index négatif du décor reste borné à elle. Son `transform`
       en crée déjà un -- mais .stage-maximized le remet à `none`, et le décor
       passerait alors derrière le fond de la scène. */
    isolation: isolate;
```

`pointer-events: none` sur les deux pseudo-éléments : ces couches couvrent toute la scène et intercepteraient sinon les clics destinés aux panneaux.

- [ ] **Step 2 : Éteindre le décor au vidéoprojecteur**

Dans le bloc `:root[data-display="projector"]` existant, ajouter :

```css
    /* Le décor s'éteint en salle mal éclairée. Ça reste dans le contrat des
       profils -- l'opacité d'une couche purement décorative ne peut déplacer
       aucune boîte -- et ça rend au texte le contraste que le profil cherche
       justement à remonter. */
    --decor-dots: 0;
    --decor-mesh: 0;
```

- [ ] **Step 3 : Régénérer et vérifier**

```bash
bash hack/css.sh engine
grep -c "stage::before\|stage::after" handlers/resources/demoit.css
grep -A3 'data-display="projector"' handlers/resources/demoit.css | grep -c "decor"
```

Attendu : les pseudo-éléments présents, et les deux jetons remis à zéro sous `projector`.

- [ ] **Step 4 : Vérifier que `sample/` reste intact**

`sample/` n'a pas de `talk.yml`, donc `Theme.Stage` est faux et son `#app` porte `.no-stage`, pas `.stage`. Les pseudo-éléments ne peuvent donc pas l'atteindre.

```bash
grep -n "no-stage" handlers/resources/index.tmpl.html
grep -c "no-stage::before" handlers/resources/demoit.css
```

Attendu : le ternaire `{{ if .Stage }}stage{{ else }}no-stage{{ end }}` toujours en place, et **zéro** occurrence de `no-stage::before`.

- [ ] **Step 5 : Commit**

```bash
git add styles/demoit.src.css handlers/resources/demoit.css
git commit -m "feat(deck): poser le décor du site sur la scène"
```

---

### Task 4 : `slide-hero`, le renforcement sur `cover` et `quote`

**Files:**
- Modify: `styles/components.css`
- Modify: `styles/demoit.src.css` (la ligne `@source inline(...)` du châssis, et le renforcement)
- Modify: `deck/layouts/cover.html`
- Modify: `deck/layouts/quote.html`
- Modify: `impact-framework/.demoit/style.css` (traitement `em`)
- Regenerate: `handlers/resources/demoit.css`, `impact-framework/.demoit/tailwind.css`

**Interfaces:**
- Consumes: `--decor-dots`, `--decor-mesh`, `--color-text-gradient`, `--color-main` (tâche 1).
- Produces: l'utilitaire `slide-hero`. La tâche 7 le pose sur la cover du talk.

- [ ] **Step 1 : Déclarer l'utilitaire**

Dans `styles/components.css`, à la suite des autres `@utility` :

```css
/*
Marque les layouts qui reçoivent le traitement complet : la cover et les
citations. Le site ne décore que son hero et laisse ses sections de contenu
plates -- même règle ici, pour que le décor ne concurrence jamais un terminal
ou du code.

L'utilitaire ne porte rien lui-même : il sert de prise à la scène, qui lit sa
présence, et au traitement `em` du talk.
*/
@utility slide-hero {
    position: relative;
}
```

- [ ] **Step 2 : L'émettre inconditionnellement**

Dans `styles/demoit.src.css`, ajouter `slide-hero` à la ligne `@source inline(...)` du châssis, qui devient :

```css
@source inline("slide-main slide-narrow slide-header slide-footer slide-prose contact-row card slide-hero");
```

Sans ça, un `@utility` n'est émis que là où un fichier scanné l'utilise — donc réordonner une layout pourrait le faire disparaître silencieusement.

- [ ] **Step 3 : Écrire le renforcement**

Dans `styles/demoit.src.css`, juste après le bloc du décor de la tâche 3 :

```css
/*
La cover et les citations poussent les deux couches.

:has() plutôt qu'une classe sur .stage : une layout ne peut pas atteindre son
propre conteneur, qui est écrit dans index.tmpl.html. Le sélecteur est supporté
par tous les navigateurs visés, le chromedp de /pdf compris.
*/
.stage:has(main.slide-hero) {
    --decor-dots: 0.10;
    --decor-mesh: 0.45;
}

.dark .stage:has(main.slide-hero) {
    --decor-mesh: 0.26;
}
```

Ces quatre valeurs sont le seul réglage de ce plan qui se décide devant le rendu : les ajuster à la tâche 8 si le décor écrase le titre ou reste invisible.

- [ ] **Step 4 : Poser la classe sur les deux layouts embarquées**

Dans `deck/layouts/cover.html`, la première ligne devient :

```
<main class="{{ if .Class }}{{ .Class }}{{ else }}slide-main slide-hero flex flex-col items-center justify-center text-center{{ end }}">
```

Dans `deck/layouts/quote.html`, de même :

```
{{ template "header" . }}<main class="{{ if .Class }}{{ .Class }}{{ else }}slide-main slide-hero slide-narrow flex flex-col items-center justify-center text-center{{ end }}">
```

**Piège à connaître :** une slide qui déclare `class:` **remplace la liste en bloc** au lieu de s'y ajouter. Une slide de cover qui déclarerait son propre `class:` perdrait donc `slide-hero` sans rien qui le signale. La tâche 7 en tient compte pour `impact-framework`.

- [ ] **Step 5 : Traiter `em` sous `slide-hero`**

Dans `impact-framework/.demoit/style.css`, juste après la règle `em` existante :

```css
/*
Sur la cover et les citations, une emphase prend le dégradé du thème et son
halo -- le pattern du hero de zatsit.fr, qui ne met en dégradé que les mots
accentués et laisse le reste du titre à contraste plein.

Ailleurs, `em` reste l'aplat d'accent défini au-dessus : sur une slide de
travail, un titre en background-clip:text ne réagirait plus aux profils
d'affichage, et c'est la ligne qu'on veut la plus lisible.
*/
.slide-hero em {
    background: var(--color-text-gradient);
    background-clip: text;
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    filter: drop-shadow(0 0 30px color-mix(in srgb, var(--color-main) 50%, transparent));
}
```

- [ ] **Step 6 : Régénérer et vérifier**

```bash
bash hack/css.sh engine
bash hack/css.sh impact-framework
grep -c "slide-hero" handlers/resources/demoit.css
GOCACHE="$TMPDIR/go-build" go test ./deck/... -run TestLayouts -v 2>&1 | tail -20
```

Attendu : `slide-hero` présent dans la feuille du moteur, et les tests de layout au vert — `TestLayoutsRenderTheEmbeddedDefault` n'assertant que sur `default`, il n'est pas touché.

- [ ] **Step 7 : Commit**

```bash
git add styles/ deck/layouts/cover.html deck/layouts/quote.html impact-framework/.demoit/style.css handlers/resources/demoit.css impact-framework/.demoit/tailwind.css
git commit -m "feat(deck): renforcer le décor sur les layouts cover et quote"
```

---

### Task 5 : Titre et sous-titre dans `talk.yml`

**Files:**
- Modify: `deck/talk.go:41-55` (la structure `Talk`) et `deck/talk.go:60-80` (`LoadTalk`)
- Test: `deck/talk_test.go`

**Interfaces:**
- Consumes: `renderTitle(title string, firstLine int) (template.HTML, int, error)`, déjà défini dans `deck/deck.go:175`, paquet `deck`, donc accessible sans export.
- Produces: `Talk.Title template.HTML` et `Talk.Subtitle template.HTML`, tous deux rendus en Markdown. La tâche 7 les lit depuis la layout de cover.

- [ ] **Step 1 : Écrire le test qui échoue**

Dans `deck/talk_test.go`, à la suite des tests existants :

```go
// Le titre et le sous-titre d'un talk passent par le même pipeline Markdown
// que le title: d'une slide, et pour la même raison : une cover écrit
// `*Impact Framework*` pour que le mot prenne le traitement d'accent, et une
// apostrophe ne doit pas ressortir en entité HTML comme le ferait
// html/template sur une chaîne nue.
func TestLoadTalkRendersTheTitleAndSubtitleAsMarkdown(t *testing.T) {
	t.Parallel()

	folder := writeTalk(t, "title: A la découverte d'*Impact Framework*.\nsubtitle: 10 avril 2025 — **Ludovic Dussart**\n")

	talk, err := deck.LoadTalk(folder)
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	if got, want := string(talk.Title), "A la découverte d'<em>Impact Framework</em>."; got != want {
		t.Errorf("got title %q, want %q", got, want)
	}
	if got, want := string(talk.Subtitle), "10 avril 2025 — <strong>Ludovic Dussart</strong>"; got != want {
		t.Errorf("got subtitle %q, want %q", got, want)
	}
}

// Un talk sans title ni subtitle rend deux chaînes vides, pas un <p></p> :
// c'est ce qui permet à une layout de les tester avec {{ with }}.
func TestLoadTalkLeavesAnAbsentTitleEmpty(t *testing.T) {
	t.Parallel()

	talk, err := deck.LoadTalk(writeTalk(t, "layout: content\n"))
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	if talk.Title != "" {
		t.Errorf("got title %q, want it empty", talk.Title)
	}
	if talk.Subtitle != "" {
		t.Errorf("got subtitle %q, want it empty", talk.Subtitle)
	}
}
```

Et adapter l'assertion existante de `TestLoadTalkReadsEveryKey`, ligne 39, qui compare un `template.HTML` à un `string` :

```go
	if got, want := string(talk.Title), "Impact Framework"; got != want {
```

- [ ] **Step 2 : Lancer le test et vérifier qu'il échoue**

```bash
GOCACHE="$TMPDIR/go-build" go test ./deck/ -run "TestLoadTalkRenders|TestLoadTalkLeaves" -v
```

Attendu : échec de compilation, `talk.Subtitle undefined` — le champ n'existe pas encore.

- [ ] **Step 3 : Écrire l'implémentation**

Dans `deck/talk.go`, ajouter `"html/template"` aux imports, puis dans la structure `Talk` remplacer le champ `Title` et ajouter `Subtitle` :

```go
	// Title is the talk's title, rendered from talk.yml through the same
	// Markdown pipeline as a slide's title: key -- so a cover can write
	// `*Impact Framework*` and have the word take the accent treatment.
	Title template.HTML `yaml:"title"`
	// Subtitle is the line a cover displays under the title, typically the
	// date and the speaker. Rendered as Markdown too.
	Subtitle template.HTML `yaml:"subtitle"`
```

`yaml.v3` désérialise sans broncher dans un type nommé de sous-jacent `string`, donc la valeur brute atterrit dans le champ avant d'être rendue.

Puis dans `LoadTalk`, entre le `yaml.Unmarshal` et le repli sur `talk.Layout` :

```go
	// Rendus après le parse, pas pendant : le champ reçoit d'abord le
	// Markdown brut, qu'on remplace par son rendu.
	title, _, err := renderTitle(string(talk.Title), 0)
	if err != nil {
		return talk, fmt.Errorf("unable to render the talk title: %w", err)
	}
	talk.Title = title

	subtitle, _, err := renderTitle(string(talk.Subtitle), 0)
	if err != nil {
		return talk, fmt.Errorf("unable to render the talk subtitle: %w", err)
	}
	talk.Subtitle = subtitle
```

`renderTitle` rend la chaîne vide telle quelle et réduit un paragraphe unique à son contenu inline, ce qui est exactement ce que les deux tests attendent.

- [ ] **Step 4 : Lancer les tests et vérifier qu'ils passent**

```bash
GOCACHE="$TMPDIR/go-build" go test ./deck/ -v 2>&1 | tail -30
```

Attendu : tous les tests du paquet `deck` au vert, `TestLoadTalkReadsEveryKey` compris.

- [ ] **Step 5 : Vérifier le reste du dépôt**

```bash
GOCACHE="$TMPDIR/go-build" go build ./... && GOCACHE="$TMPDIR/go-build" go test ./... 2>&1 | tail -12
gofmt -l deck/
```

Attendu : build propre, tous les paquets au vert, et `gofmt -l` silencieux.

- [ ] **Step 6 : Commit**

```bash
git add deck/talk.go deck/talk_test.go
git commit -m "feat(deck): rendre le titre et le sous-titre d'un talk en Markdown"
```

---

### Task 6 : Logos partenaires clair et sombre

**Files:**
- Modify: `deck/talk.go` (structure `Talk`)
- Modify: `deck/layouts/partials.html` (le bloc `header`)
- Modify: `deck/layouts/cover.html`
- Test: `deck/talk_test.go`, `deck/layout_test.go`

**Interfaces:**
- Consumes: `Talk.Logos []string`, déjà en place.
- Produces: `Talk.LogosDark []string`. La tâche 7 le renseigne dans `talk.yml`.

`logo.jpg` est un JPEG à fond blanc opaque : sur une slide sombre c'est un rectangle blanc. On ne détoure rien ici — on pose la mécanique, les fichiers viendront.

- [ ] **Step 1 : Écrire les tests qui échouent**

Dans `deck/talk_test.go` :

```go
func TestLoadTalkReadsTheDarkLogos(t *testing.T) {
	t.Parallel()

	folder := writeTalk(t, "logos:\n  - /images/a.svg\nlogosDark:\n  - /images/a-dark.svg\n")

	talk, err := deck.LoadTalk(folder)
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	if got, want := len(talk.LogosDark), 1; got != want {
		t.Fatalf("got %d dark logos, want %d", got, want)
	}
	if got, want := talk.LogosDark[0], "/images/a-dark.svg"; got != want {
		t.Errorf("got dark logo %q, want %q", got, want)
	}
}
```

Dans `deck/layout_test.go` :

```go
// Sans logosDark, aucune classe conditionnelle n'est émise : un talk qui ne
// déclare pas de variante sombre doit rendre exactement comme avant, dans les
// deux thèmes. C'est ce qui protège sample/ et tout talk existant.
func TestHeaderLogosCarryNoThemeClassWithoutDarkLogos(t *testing.T) {
	t.Parallel()

	slide := deck.Slide{Talk: deck.Talk{Logos: []string{"/images/a.svg"}}}

	got, err := deck.NewLayouts(t.TempDir()).Execute("default", slide)
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	if !strings.Contains(string(got), `src="/images/a.svg"`) {
		t.Errorf("got %q, want it to contain the logo", got)
	}
	if strings.Contains(string(got), "dark:hidden") {
		t.Errorf("got %q, want no dark:hidden when the talk declares no dark logos", got)
	}
}

// Avec logosDark, les deux jeux sont émis et la bascule est en CSS : rien côté
// Go ne connaît le thème courant, qui n'est décidé que dans le navigateur.
func TestHeaderEmitsBothLogoSetsWithDarkLogos(t *testing.T) {
	t.Parallel()

	slide := deck.Slide{Talk: deck.Talk{
		Logos:     []string{"/images/a.svg"},
		LogosDark: []string{"/images/a-dark.svg"},
	}}

	got, err := deck.NewLayouts(t.TempDir()).Execute("default", slide)
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	for _, want := range []string{
		`src="/images/a.svg"`,
		"dark:hidden",
		`src="/images/a-dark.svg"`,
		"hidden dark:block",
	} {
		if !strings.Contains(string(got), want) {
			t.Errorf("got %q, want it to contain %q", got, want)
		}
	}
}
```

- [ ] **Step 2 : Lancer les tests et vérifier qu'ils échouent**

```bash
GOCACHE="$TMPDIR/go-build" go test ./deck/ -run "DarkLogos|LogoSets|ThemeClass" -v
```

Attendu : échec de compilation, `LogosDark undefined`.

- [ ] **Step 3 : Ajouter le champ**

Dans `deck/talk.go`, après `Logos` :

```go
	// LogosDark are the dark-theme counterparts of Logos, in the same order.
	// A talk that declares none renders Logos in both themes, which is what
	// every talk did before this field existed.
	LogosDark []string `yaml:"logosDark"`
```

- [ ] **Step 4 : Émettre les deux jeux dans le header**

Dans `deck/layouts/partials.html`, remplacer la boucle des logos du bloc `header` par :

```
    <div class="flex items-center gap-4">
        {{ range .Talk.Logos }}<img class="size-12 rounded-full object-cover{{ if $.Talk.LogosDark }} dark:hidden{{ end }}" src="{{ . }}">
        {{ end }}{{ range .Talk.LogosDark }}<img class="size-12 rounded-full object-cover hidden dark:block" src="{{ . }}">
        {{ end }}
    </div>
```

Le `{{ if $.Talk.LogosDark }}` est le cœur du mécanisme : sans variante sombre, aucune classe conditionnelle n'est émise et le rendu est identique à aujourd'hui. `$` désigne le contexte racine, la `Slide`, la variable `.` étant réaffectée à l'élément courant à l'intérieur d'un `range`.

- [ ] **Step 5 : Faire de même sur la cover embarquée**

Dans `deck/layouts/cover.html`, remplacer la boucle par :

```
<div class="mt-12 flex items-center justify-center gap-12">
    {{ range .Talk.Logos }}<img class="h-[25rem] w-[25rem] object-contain{{ if $.Talk.LogosDark }} dark:hidden{{ end }}" src="{{ . }}" loading="lazy">
    {{ end }}{{ range .Talk.LogosDark }}<img class="h-[25rem] w-[25rem] object-contain hidden dark:block" src="{{ . }}" loading="lazy">
    {{ end }}
</div>
```

- [ ] **Step 6 : Lancer les tests et vérifier qu'ils passent**

```bash
GOCACHE="$TMPDIR/go-build" go test ./deck/... -v 2>&1 | tail -30
```

Attendu : tout le paquet au vert, `TestLayoutsRenderTheEmbeddedDefault` compris — il n'assertait que sur `src="/images/a.svg"`, qui ne bouge pas.

- [ ] **Step 7 : Régénérer la feuille du moteur**

`hidden`, `dark:hidden` et `dark:block` sont des utilitaires Tailwind standards, vus par le scanner dans `deck/layouts/`. Ils n'existent pas tant que le build n'a pas tourné.

```bash
bash hack/css.sh engine
grep -c "dark\\\\:hidden\|dark\\\\:block" handlers/resources/demoit.css
```

Attendu : au moins une occurrence de chaque.

- [ ] **Step 8 : Commit**

```bash
git add deck/talk.go deck/talk_test.go deck/layout_test.go deck/layouts/partials.html deck/layouts/cover.html handlers/resources/demoit.css
git commit -m "feat(deck): accepter une variante sombre des logos d'un talk"
```

---

### Task 7 : La cover d'`impact-framework`

**Files:**
- Create: `impact-framework/.demoit/images/logo-full.svg`
- Modify: `impact-framework/.demoit/layouts/cover.html`
- Modify: `impact-framework/.demoit/talk.yml`
- Modify: `impact-framework/.demoit/style.css`
- Modify: `impact-framework/demoit.md:1-8`
- Regenerate: `impact-framework/.demoit/tailwind.css`

**Interfaces:**
- Consumes: `Talk.Title`, `Talk.Subtitle` (tâche 5), `Talk.LogosDark` (tâche 6), `slide-hero` (tâche 4), `--logo-gradient-*` et `--color-main` (tâche 1).

- [ ] **Step 1 : Copier l'asset**

```bash
cp ../../../../Documents/Workplace/zatsit/zats-websites/corporate/src/assets/icons/logo-full.svg \
   impact-framework/.demoit/images/logo-full.svg
```

Si le chemin relatif ne résout pas depuis le worktree, utiliser le chemin absolu `/Users/ex5285/Documents/Workplace/zatsit/zats-websites/corporate/src/assets/icons/logo-full.svg`.

Vérifier que le fichier porte bien `fill="currentColor"` sur la racine et la def `linearGradient id="logo-full-gradient"` dont les stops lisent `--logo-gradient-start/mid/end` :

```bash
grep -c "logo-full-gradient\|currentColor" impact-framework/.demoit/images/logo-full.svg
```

Attendu : au moins 2.

- [ ] **Step 2 : Renseigner `talk.yml`**

`impact-framework/.demoit/talk.yml` devient :

```yaml
layout: default
title: A la découverte d'*Impact Framework*.
subtitle: 10 avril 2025 — Ludovic Dussart
theme:
  # La scène fixe 1920x1080, mise à l'échelle de ce sur quoi on projette.
  stage: true
  # Rend le commutateur de thème : t bascule clair/sombre, d fait défiler le
  # profil d'affichage (screen / tv / projector).
  dark: true
logos:
  - /images/zatsit_logo.svg
  - /images/logo.jpg
# Variantes sombres, dans le même ordre que logos. Tant que cette liste est
# absente ou vide, les logos ci-dessus servent dans les deux thèmes.
# logosDark:
#   - /images/zatsit_logo_blanc.svg
#   - /images/logo-dark.png
```

La liste `logosDark` reste commentée : la mécanique est prête, les fichiers viendront. `zatsit_logo_blanc.svg` existe déjà et ferait un premier candidat ; `logo.jpg` attend une version à fond transparent.

- [ ] **Step 3 : Réécrire la layout de cover du talk**

`impact-framework/.demoit/layouts/cover.html` en entier :

```
{{/* Le SVG est inliné, pas chargé en <img>, et c'est nécessaire : la def
     #logo-full-gradient et la référence url(#logo-full-gradient) doivent
     vivre dans le même document. C'est ce que fait Hero.astro sur zatsit.fr
     via <LogoFull />. Un <img> ne peut pas être recoloré par le thème, et un
     <use> vers un <symbol> ne marcherait pas non plus : un sélecteur du
     document n'atteint pas l'arbre fantôme d'un <use>, et les path
     retomberaient sur le noir par défaut.

     slide-hero est écrit ici plutôt que laissé au défaut de la layout
     embarquée, que ce fichier remplace : un talk qui surcharge une layout
     possède aussi ses classes par défaut. */}}
<main class="{{ if .Class }}{{ .Class }}{{ else }}slide-main slide-hero flex flex-col items-center justify-center text-center{{ end }}">

    <svg class="logo-full-gradient-glow w-[40rem] max-w-[70%] h-auto" viewBox="0 0 523.94 426.01" xmlns="http://www.w3.org/2000/svg" role="img" aria-label="zatsit">
        <defs>
            <linearGradient id="logo-full-gradient" gradientUnits="userSpaceOnUse" x1="0" y1="0" x2="523.94" y2="426.01">
                <stop offset="0%" style="stop-color: var(--logo-gradient-start, #0F15FD)"/>
                <stop offset="50%" style="stop-color: var(--logo-gradient-mid, #0a0ecc)"/>
                <stop offset="100%" style="stop-color: var(--logo-gradient-end, #020466)"/>
            </linearGradient>
        </defs>
        <path d="m301.52,171.11c-41.58.67-76.15,27.12-88.31,65.51,14.31,5.83,29.95,9.07,46.35,9.07,50.77,0,94.33-30.8,113.06-74.72-23.7-.1-47.41-.25-71.1.13"/>
        <path d="m220.88,74.73c42.24-.48,77.5-26.84,89.88-63.55-15.59-7.16-32.92-11.18-51.19-11.18-50.8,0-94.39,30.83-113.09,74.8,24.8.06,49.6.22,74.4-.06"/>
        <path d="m350.17,144.03c13.76-8.48,24.38-18.76,31.85-30.68-3.06-40.01-25.27-74.64-57.5-94.76-21.02,43.1-42.11,86.34-63.45,130.09,30.83,13.61,60.82,12.78,89.1-4.65"/>
        <path d="m141.61,135.88c-1.21,2.18-2.29,4.4-3.3,6.64,6.04,37.5,29.05,69.31,60.9,87.32,21.65-44.38,43.18-88.53,64.69-132.63-30.73-18.56-93.04-13.98-122.29,38.67"/>
        <path d="m425.11,327.29h-18.5v96.19h18.5v-96.19ZM0,423.47h88.6v-16.55H24.73l61.92-64.45v-15.19H2.92v16.55h59L0,408.29v15.19Zm295.16-27.6h-18.84c-2.65,9.01-9.63,12.99-19.04,12.99-12.07,0-20.06-6.62-20.06-21.81v-42.84h38.97v-16.94h-38.97v-31.07h-18.5v90.84c0,26.09,15.58,38.94,38.55,38.94,20.22,0,34.69-9.96,37.88-30.13m89.85-41.53c-.39-16.55-14.6-29.6-37.58-29.6s-37.77,12.85-37.77,31.74c0,22,20.25,25.12,36.99,27.26,12.07,1.56,22.2,2.73,22.2,11.88,0,7.59-7.21,13.63-21.22,13.63s-20.83-5.65-21.22-13.63h-18.69c.58,17.53,15.77,30.38,39.92,30.38s39.92-12.66,39.92-32.13c0-22-20.44-25.12-37.19-27.26-11.88-1.56-22-2.53-22-11.68,0-7.59,7.01-13.44,19.28-13.44s18.5,5.45,19.08,12.85h18.3Z"/>
        <path d="m523.94,395.88h-18.84c-2.65,9.02-9.63,12.99-19.04,12.99-12.07,0-20.06-6.62-20.06-21.81v-42.84h38.97v-16.94h-38.97v-31.06h-18.5v90.84c0,26.09,15.58,38.94,38.55,38.94,20.22,0,34.69-9.96,37.88-30.13"/>
        <path d="m182.8,372.99c0,18.04-14.81,32.66-33.07,32.66s-33.07-14.62-33.07-32.66,14.81-32.66,33.07-32.66,33.07,14.62,33.07,32.66m17.93,50.49v-50.49h-.04c-.47-27.02-23.09-48.79-50.96-48.79s-51,22.22-51,49.64,22.83,49.64,51,49.64c12.63,0,24.17-4.49,33.07-11.9v11.9h17.93Z"/>
    </svg>

    {{ with .Talk.Title }}<h1 class="mt-12 mb-0 text-6xl font-bold leading-tight tracking-tight">{{ . }}</h1>
    {{ end }}{{ with .Talk.Subtitle }}<p class="mt-6 text-3xl font-medium text-fg-muted">{{ . }}</p>
    {{ end }}

    {{/* Le contenu de la slide reste rendu : une cover peut vouloir ajouter
         une ligne que talk.yml n'a pas à porter. Vide aujourd'hui. */}}
    {{ .Content }}

    <div class="mt-16 flex items-center justify-center gap-16">
        {{ range .Talk.Logos }}<img class="h-24 w-auto max-w-[16rem] object-contain{{ if $.Talk.LogosDark }} dark:hidden{{ end }}" src="{{ . }}" loading="lazy">
        {{ end }}{{ range .Talk.LogosDark }}<img class="h-24 w-auto max-w-[16rem] object-contain hidden dark:block" src="{{ . }}" loading="lazy">
        {{ end }}
    </div>
</main>
{{ template "notes" . }}
```

Les logos passent de `25rem` carrés à `h-24` (6rem) de haut : le lockup zatsit occupe désormais le haut de la slide, les logos redeviennent une mention de pied. `w-auto` plutôt qu'un carré, parce qu'`object-contain` dans une boîte carrée laissait de larges vides de part et d'autre d'un logo large.

- [ ] **Step 4 : Reprendre les deux règles du site**

Dans `impact-framework/.demoit/style.css`, après le traitement `em` de la tâche 4 :

```css
/*
Les deux règles de zatsit.fr, reprises mot pour mot. Le sélecteur vise `path`
et non l'élément : le fill de la racine du SVG est `currentColor`, et c'est la
def #logo-full-gradient qui doit gagner sur chaque tracé.

Correct uniquement parce que le SVG est un vrai élément du document. Avec un
<use> vers un <symbol>, ce sélecteur n'atteindrait pas l'arbre fantôme et les
path retomberaient sur le noir.
*/
.logo-full-gradient-glow path {
    fill: url(#logo-full-gradient);
}

.logo-full-gradient-glow {
    filter: drop-shadow(0 0 40px color-mix(in srgb, var(--color-main) 40%, transparent));
}
```

Et **supprimer** la plaque `.title`, devenue sans objet :

```css
/* The cover's gradient plate. */
.title {
    background: linear-gradient(to bottom, var(--color-main) 10%, var(--color-surface) 90%);
}
```

- [ ] **Step 5 : Nettoyer la slide de cover**

Dans `impact-framework/demoit.md`, les huit premières lignes deviennent :

```
---
layout: cover
---

---
```

Disparaissent : la ligne `class:` — qui écraserait `slide-hero`, puisqu'un `class:` remplace la liste par défaut en bloc — et les deux `<h1>` / `<h3>` forcés en `style="color: white"`, un blanc qui n'a plus de plaque bleue sous lui. Le titre vient maintenant de `talk.yml`.

Ce sont les deux seules couleurs en dur du deck ; les autres occurrences sont dans le `classDef` mermaid, qui garde les siennes.

- [ ] **Step 6 : Régénérer et vérifier**

```bash
bash hack/css.sh impact-framework
grep -c "logo-full-gradient" impact-framework/.demoit/layouts/cover.html
grep -c "style=\"color: white\"" impact-framework/demoit.md
GOCACHE="$TMPDIR/go-build" go test ./... 2>&1 | tail -12
```

Attendu : la def présente dans la layout, **zéro** couleur blanche en dur dans le deck, tous les paquets au vert.

- [ ] **Step 7 : Commit**

```bash
git add impact-framework/
git commit -m "feat(impact-framework): refondre la cover autour du logo zatsit en dégradé"
```

---

### Task 8 : Revue visuelle et mise à jour du CLAUDE.md

**Files:**
- Modify: `CLAUDE.md`
- Possiblement: `styles/demoit.src.css` (ajustement des quatre valeurs de décor renforcé)

C'est la seule tâche qui ne peut pas se vérifier en ligne de commande. La suite Go ne dit rien de ce que le navigateur peint, et ce plan ne change presque que ça.

- [ ] **Step 1 : Lancer le talk**

```bash
GOCACHE="$TMPDIR/go-build" go build -o demoit . && ./demoit --dev impact-framework
```

Ouvrir `http://localhost:8888`.

- [ ] **Step 2 : Parcourir la checklist**

Sur chaque point, basculer le thème avec `t` et faire défiler les profils avec `d` :

- **la cover** — le lockup prend-il bien le dégradé (bleu → marine en clair, or → ambre → rouge en sombre) et non du noir ? Un logo noir signerait un sélecteur qui n'atteint pas les `path` ;
- **le titre de la cover** — `Impact Framework` en dégradé, le reste à contraste plein ;
- **le décor** — visible sans écraser le texte sur les slides de travail, plus présent sur la cover. C'est ici qu'on ajuste les quatre valeurs de la tâche 4, step 3 ;
- **le profil `projector`** — le décor disparaît complètement ;
- **la figure mermaid en sombre** (slide « IF - Concepts majeurs ») — ses libellés sont déjà rognés à droite, défaut connu et documenté : vérifier que le décor ne l'aggrave pas, et surtout que le fond blanc de `.dark .mermaid` tient toujours ;
- **une slide `split` avec un tty et une iframe** (slide « IF - Hello World ! ») — rien ne passe au-dessus des panneaux, aucun scintillement, le terminal reste utilisable. Tester aussi le maximisé, le zoom d'un panneau, et le retour ;
- **`/grid`** et **`/speakernotes`** — suivent le thème par `BroadcastChannel` ;
- **`/pdf`** — force `?theme=light`, donc sort en clair avec le décor figé.

- [ ] **Step 3 : Vérifier que `sample/` est intact**

```bash
./demoit sample
```

Attendu : rendu identique à avant ce chantier — pas de scène, pas de commutateur de thème, pas de décor. `sample/` n'a pas de `talk.yml`.

- [ ] **Step 4 : Corriger le CLAUDE.md**

Deux affirmations y sont fausses ou périmées.

La première, dans la section « Conventions », dit que la palette d'un talk vit dans son `tokens.css`. C'est vrai de la génération des utilitaires et faux du rendu. Remplacer la phrase par :

```
- Palette et per-talk theming : les valeurs qui s'appliquent réellement vivent
  dans le moteur — `styles/tokens.css` pour le clair, le bloc `.dark` de
  `styles/demoit.src.css` pour le sombre. `<folder>/.demoit/tokens.css` en est
  le **reflet** : importé avec `theme(reference)`, il donne le vocabulaire au
  générateur d'utilitaires et **n'émet aucune variable** — `tailwind.css` ne
  contient aucun bloc `:root`. Faire diverger les deux ne repeint rien. Les
  passer en import simple ne marche pas non plus : le `:root, :host` alors émis
  est hors couche, de même spécificité que le `.dark` du moteur, et
  `/tailwind.css` chargeant après `/demoit.css`, le clair gagnerait sur le
  sombre. **Ce fichier doit contenir des blocs `@theme` et rien d'autre.**
```

La seconde : la section sur `talk.yml` dit qu'il « parse `title` et `footer`, mais aucun layout ne lit l'un ou l'autre ». `title` est maintenant lu par la cover d'`impact-framework`, et `subtitle` et `logosDark` s'ajoutent. Mettre la phrase à jour, et documenter le décor de la scène et `slide-hero` dans la section « The rendering chassis ».

- [ ] **Step 5 : Vérification finale**

```bash
GOCACHE="$TMPDIR/go-build" go build ./... && GOCACHE="$TMPDIR/go-build" go test ./...
gofmt -l .
bash hack/css.sh engine && bash hack/css.sh impact-framework
git status --short
```

Attendu : build et tests au vert, `gofmt -l` silencieux, et **`git status` vide après les rebuilds** — c'est ce qui prouve que les artefacts commités correspondent bien à leurs sources.

- [ ] **Step 6 : Commit**

```bash
git add CLAUDE.md styles/ handlers/resources/demoit.css
git commit -m "docs: corriger ce que CLAUDE.md dit de l'émission des jetons"
```
