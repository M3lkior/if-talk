# Migration beercss → Tailwind, scène fixe, thème sombre

Date : 2026-09-10
Statut : design validé, prêt pour plan d'implémentation

## Problème

Trois symptômes, une seule cause.

**Poser une classe au bon endroit relève de la devinette.** `max` veut dire deux choses différentes dans le même repo : `max-inline-size:100%` sur `<main class="responsive max">` (`deck/layouts/bare.html:1`), mais `flex:1` sur `<h2 class="max center-left">` parce que la règle `:is(nav,.row)>.max{flex:1}` gagne dans un `<nav>` (`deck/layouts/partials.html:3`). `responsive` vaut `inline-size:-webkit-fill-available` sur un bloc, mais `block-size:3rem; inline-size:3rem; object-fit:cover` sur une `<img>` — donc les logos du header partial sont des vignettes de 3 rem recadrées, ce qui n'est visiblement pas l'intention. Et `center-left`, `center-right`, `padding`, utilisées 11 fois dans les layouts, **n'existent pas dans beercss** : zéro règle, aucun effet, depuis toujours.

**Le rendu change d'un écran à l'autre.** `impact-framework/.demoit/style.css:28` fait de `rem` un `3vw`, donc les `3.5625rem` de `h1` valent 205 px sur un 1920 et 137 px sur un 1280 — pendant que la grille beercss bascule ses colonnes sur des seuils de **largeur** (601 px, 993 px). Une slide dessinée sur un écran 16:10 déporté et projetée en 16:9 change à la fois de taille de police et de nombre de colonnes. Les hauteurs nommées aggravent le tableau : `.xlarge-height{block-size:64rem}` vaut 3686 px à 57.6 px/rem, soit 3,4 fois la hauteur utile.

**Il n'y a pas de thème sombre.** `<body class="light">` est figé dans `handlers/resources/index.tmpl.html:19`, et les composants web portent leurs couleurs en dur dans leurs styles de shadow DOM.

Objectif : remplacer beercss par Tailwind CSS, rendre la géométrie déterministe quel que soit l'écran de projection, et introduire un thème sombre commutable que les runtimes embarqués suivent.

## Décisions

| # | Décision | Motif |
|---|---|---|
| 1 | Tailwind v4 via le **CLI standalone épinglé**, CSS généré et commité | Aucun `node_modules` dans un repo Go à deps vendorées ; surtout, **aucune dépendance réseau pendant un talk** — le réseau ne sert qu'au build |
| 2 | **Scène fixe 1920×1080 mise à l'échelle**, pas de breakpoints | `/pdf` (1920×1080@2x) et `/grid` (iframe 1920×1080) supposent déjà cette taille. Dessiner une fois, projeter partout ; et l'iframe du tty garde une largeur logique constante, donc la démo ne se reflowe plus selon la salle |
| 3 | **Profils d'affichage** `screen` / `tv` / `projector` en plus du thème | La scène règle la géométrie, pas la lisibilité : un rétroprojecteur délavé exige plus de contraste et de graisse, pas une autre mise en page |
| 4 | **`transform: scale()`**, pas de mise à l'échelle par `rem` | Le scaling par `rem` donnerait à l'iframe du tty sa taille physique, donc gotty recalculerait ses colonnes selon la salle — exactement le reflow que la décision 2 supprime |
| 5 | Correspondance **hybride** : utilitaires Tailwind standards + quelques `@utility` nommés | Grille, alignements, espacements et tailles ont un équivalent direct. Les 3-4 motifs récurrents du deck (carte, ligne de contact, spacer de scène) restent une classe lisible plutôt que douze utilitaires |
| 6 | Scène et thème **opt-in par talk** (bloc `theme:` dans `talk.yml`) | Protège `sample/` de la scène fixe et d'un commutateur dont son `demoit.js` ne suivrait pas la chrome. Ce que l'opt-in **ne peut pas** protéger : voir « Le sort de `sample/` » |
| 7 | Périmètre : `impact-framework/demoit.md` **et** `demoit-en.html` | `index.tmpl.html` est embarqué et partagé : retirer beercss casse tout deck qui s'y appuyait, et `demoit-en.html` s'y appuie sur 584 lignes |
| 8 | **Socle seulement**, pas de refonte esthétique | Les jetons de design posés ici rendent l'embellissement ultérieur trivial et réversible ; c'est la main du présentateur, pas celle du moteur |

## Architecture CSS

### Trois fichiers, un rôle chacun

La contrainte qui impose ce découpage : `handlers/resources/*.tmpl.html` et `deck/layouts/*.html` sont **embarqués dans le binaire** (`//go:embed`). Leurs classes ne peuvent donc pas être scannées depuis un dossier de talk — qui installe `demoit` et crée un talk n'a ni `deck/layouts/` ni `handlers/resources/` sous la main.

| Fichier | Source de vérité | Construit par | Contenu |
|---|---|---|---|
| `handlers/resources/demoit.css` | `styles/demoit.src.css` | CLI Tailwind, scan de `handlers/resources/*.tmpl.html` + `deck/layouts/*.html` | preflight, jetons `@theme`, classes sémantiques du moteur, machinerie scène / thème / profils, custom properties de chrome |
| `<talk>/.demoit/tailwind.css` | `<talk>/.demoit/tailwind.src.css` | CLI Tailwind, scan de `<talk>/**` | les utilitaires que les slides du talk utilisent, **sans** preflight ni jetons dupliqués |
| `<talk>/.demoit/style.css` | — | à la main | rôle inchangé : palette du talk (`--color-main`), débarrassée de ce qui devient un jeton |

`demoit.css` est embarqué et servi sur `/demoit.css` ; `tailwind.css` est servi par `handlers.Static` sur une nouvelle route `/tailwind.css`, comme `/style.css`. Les deux sont cache-bustés par la fonction de template `hash` déjà en place.

Séparation des couches Tailwind v4 pour éviter de livrer deux fois le reset et les jetons :

```css
/* styles/demoit.src.css — le moteur */
@layer theme, base, components, utilities;
@import "tailwindcss/theme.css" layer(theme);
@import "tailwindcss/preflight.css" layer(base);
@import "tailwindcss/utilities.css" layer(utilities);
@import "./tokens.css";

/* <talk>/.demoit/tailwind.src.css — le talk */
@import "tailwindcss/utilities.css" layer(utilities);
@import "./tokens.css" theme(reference);   /* n'émet rien, donne le vocabulaire */
```

`tokens.css` est le fichier de jetons partagé. Il est copié dans `<talk>/.demoit/`, comme `demoit.js` et `style.css` le sont déjà (« Adding a talk = copy an existing `.demoit/js/demoit.js` and `style.css` as the starting point », CLAUDE.md).

### La chaîne de build

Binaire standalone Tailwind **v4.3.3** (`tailwindcss-macos-arm64`, 76 Mo), téléchargé par `hack/css.sh` dans `.tools/` — gitignoré, version épinglée dans le script. Le réseau ne sert qu'à cette étape.

```bash
hack/css.sh engine                    # styles/demoit.src.css -> handlers/resources/demoit.css
hack/css.sh impact-framework          # .demoit/tailwind.src.css -> .demoit/tailwind.css
hack/css.sh impact-framework --watch  # à utiliser avec demoit --dev
```

Le livereload de `--dev` surveille déjà le dossier du talk, donc un `tailwind.css` régénéré déclenche le rafraîchissement sans travail supplémentaire.

### Les classes que Tailwind ne peut pas voir

`deck/directive/transform.go:107` construit ses classes **en Go** : le scanner ne les rencontrera jamais dans un fichier source. L'entrée du moteur les déclare explicitement :

```css
@source inline("col-span-{1,2,3,4,5,6,7,8,9,10,11,12}");
@source inline("h-stage-{sm,md,lg,xl}");
```

Oublier cette déclaration produit exactement le mode de panne que `transform.go` documente déjà pour les poids hors 1..12 : une classe qu'aucune feuille ne définit, un panneau qui prend toute la ligne sur scène, et rien qui dise pourquoi.

### Ce qui disparaît de `index.tmpl.html`

Chaque suppression a été vérifiée, pas supposée :

| Ligne | Élément | Vérification |
|---|---|---|
| 7 | `<link>` beercss | remplacé |
| 42 | `beer.min.js` | aucun appel `ui()` ni autre API beercss dans les deux `demoit.js` |
| 43 | `material-dynamic-colors.min.js` | référencé nulle part |
| 21-23 | `<dialog id="maximized">` | `#maximized` référencé nulle part ; le bouton vert des fenêtres agit sur son propre shadow DOM (`demoit.js:181`) |
| 19 | `class="light"` sur `<body>` | remplacé par la mécanique de thème |

`.demoit/poppins.css` (107 lignes) est également supprimé : ses classes sont déjà dupliquées dans `style.css:236-342`, et aucun template ne le charge.

Poppins passe du CDN Google à `<talk>/.demoit/fonts/`, servi par la route `/fonts/` qui existe et ne sert rien aujourd'hui. Quatre graisses seulement — 100, 400, 500, 700 — les seules utilisées, la 500 servant au profil `projector`. mermaid reste importé du CDN par `demoit.js:630` : c'est le dernier point de dépendance réseau à l'exécution, hors périmètre ici.

## La scène

`#app` devient `.stage` :

```css
.stage {
  position: absolute;
  inset-block-start: 50%;
  inset-inline-start: 50%;
  inline-size: 1920px;
  block-size: 1080px;
  transform: translate(-50%, -50%) scale(var(--stage-scale, 1));
}
```

`--stage-scale` vaut `min(innerWidth / 1920, innerHeight / 1080)`, posé par un listener de `resize` de quelques lignes dans `demoit.js`. La division longueur / longueur en CSS pur (`calc(100vw / 1920px)`) est spécifiée mais pas encore implémentée dans les moteurs, d'où le JS.

**`font-size` de la racine passe de `3vw` à `57.6px` fixes** — exactement `3vw` de 1920. Conséquence : toutes les tailles en `rem` du deck valent ce qu'elles valent aujourd'hui sur un écran 1920, puis la scène entière est mise à l'échelle. Sur n'importe quel viewport 16:9, **le rendu est identique à l'actuel au pixel près**. Seul un viewport non-16:9 diffère : il letterboxe au lieu de s'étirer, ce qui est l'objectif.

Le fond de letterbox est un jeton (`--stage-void`), aujourd'hui noir de fait (`style.css:27`).

## Les profils d'affichage

`data-display` sur `<html>`, valeurs `screen` (défaut), `tv`, `projector`. **Règle structurante : un profil ne redéfinit que des jetons de lisibilité, jamais de géométrie.** C'est ce qui garantit qu'un profil ne peut pas casser un layout.

| Jeton | `screen` | `tv` | `projector` |
|---|---|---|---|
| `--weight-body` | 400 | 500 | 500 |
| `--weight-thin` | 100 | 300 | 400 |
| `--rule` | 1px | 1.5px | 2px |
| `--rule-opacity` | .4 | .7 | 1 |
| `--fg` | `#3d4043` | `#2a2d30` | `#000` |
| `--accent` | `var(--color-main)` | `var(--color-main)` | variante assombrie |

Cibles déjà repérées dans le deck : le footer en `poppins-thin` (graisse 100, invisible sur un rétro), `#progression{opacity:.4}` (`style.css:134`), le `border-bottom: solid 1px` du header (`style.css:87`), et le `em{color:var(--color-main)}` (`style.css:124`) qui bave sur un rétro délavé.

## Le thème sombre

```css
@custom-variant dark (&:where(.dark, .dark *));
```

Variante pilotée par une classe, pas par `prefers-color-scheme` : le présentateur commute, la machine ne décide pas.

**Pas de flash.** Un script inline dans le `<head>` pose `class="dark"` et `data-display` sur `documentElement` depuis `localStorage` avant le premier paint. Sans lui, chaque navigation de slide clignoterait en blanc.

**Opt-in.** `deck/talk.go` gagne un bloc `theme:` à deux booléens, et `handlers/step.go` les expose au template — `deck.LoadTalk` est déjà exporté :

```yaml
# <talk>/.demoit/talk.yml
theme:
  stage: true   # scène fixe 1920x1080 ; sans elle, le talk reste pleine fenêtre
  dark: true    # commutateur de thème ; sans elle, aucun commutateur n'est rendu
```

`index.tmpl.html` conditionne l'enveloppe `.stage` et le script de thème sur ces deux clés. `impact-framework` les active, `sample` ne les déclare pas.

**Commutation et propagation.** Touche `t` pour le thème, `d` pour cycler les profils, l'une et l'autre persistées dans `localStorage` et rediffusées sur le `BroadcastChannel("demoit_nav")` déjà utilisé par les notes. Aucun contrôle à l'écran : la scène reste propre, et `<nav-arrows>` ne réserve que les flèches, `PageUp`/`PageDown` et l'espace (`demoit.js:519-533`), donc `t` et `d` sont libres. `/speakernotes` et `/grid` sont sur la même origine, ils lisent le même `localStorage` et écoutent le même canal. Un paramètre d'URL `?theme=light` sert d'override — `/pdf` s'en sert pour rendre en clair, sur le précédent du `?grid=true` existant.

**La chrome des composants passe par des custom properties**, seul mécanisme qui traverse le shadow DOM (les classes, non — c'est pourquoi `class="grid center-right s6 circle transparent"` dans `TitleBar.render()` et `class="large-height"` dans `SourceCode.render()` sont morts aujourd'hui, et seront supprimés). Les couleurs en dur de `impact-framework/.demoit/js/demoit.js` deviennent des `var(--dm-*)` :

| Aujourd'hui | Où | Jeton |
|---|---|---|
| `#ddd`, `white` | `FakeWindow` `.main` | `--dm-window-bg` |
| `linear-gradient(#edeaed, #dddfdd)`, `#cbcbcb` | `FakeWindow` `#bar` | `--dm-chrome-bg`, `--dm-chrome-border` |
| `rgb(243,243,243)`, `rgb(236,236,236)` | `SourceCode` `#tabs` | `--dm-tabs-bg`, `--dm-tab-bg` |
| `#212121` | `SourceCode` `.chroma` | `--dm-code-fg` |
| `rgb(191,214,255)` | `SourceCode` `--default-color-selection` | `--dm-code-selection` |

`sample/.demoit/js/demoit.js` n'est pas touché.

**Les runtimes embarqués.**

| Runtime | Suit le sombre | Comment |
|---|---|---|
| `source-code` | oui | re-fetch de `/sourceCode/...&style=github-dark` au basculement. chroma émet déjà sa feuille avec le HTML (`html.Standalone(true)`, `handlers/code.go:50`), donc rien à écrire côté CSS et `handlers/code.go` ne change pas |
| mermaid | oui | `mermaid.initialize({theme})` puis re-run après reset de `data-processed` |
| notes, `/grid` | oui | `localStorage` + `BroadcastChannel` |
| chrome des fenêtres | oui | custom properties ci-dessus |
| `web-term` (gotty) | **non, sombre en permanence** | gotty ne fixe son thème qu'au démarrage du serveur (`shell/shell.go:17`) ; le commuter exigerait un redémarrage, et un tty clair n'est attendu par personne |
| `/pdf` | **non, clair en permanence** | c'est de l'impression |
| `vs-code` | hors périmètre | exigerait de monter un user-data-dir généré dans le conteneur (`vscode/server.go:94` ne monte que le dossier du talk) plus une route d'écriture du `settings.json` |
| `web-browser` | hors périmètre | iframe vers une application tierce ; seul un paramètre d'URL qu'elle accepte pourrait agir |

L'accent `--color-main: #0f15fd` est quasi noir sur fond sombre : sa variante sombre est un bleu de même teinte remonté en luminosité, calibré à l'œil sur les premières slides.

## Correspondance des classes

Établie depuis les règles réelles de beercss 3.7.8, pas depuis sa documentation. Seules les classes que le repo utilise.

| beercss | effet réel | remplacement |
|---|---|---|
| `.grid` | `repeat(12, calc(8.33% - gap + gap/12))`, `gap:1rem` | `grid grid-cols-12 gap-4` |
| `.s1`…`.s12` | `grid-area:auto/span N` | `col-span-N` |
| `.m6` / `.l6` | span 6 au-delà de 601 / 993 px | **supprimés** : une scène fixe n'a qu'une largeur. `s12 m6 l6` → `col-span-6` |
| `.large-space` / `.medium-space` / `.small-space` | div de `block-size` 3 / 2 / 1 rem | `h-12` / `h-8` / `h-4` |
| `.center-align` | `text-align:center` + `justify-content:center` | `text-center justify-center` |
| `.left-align` / `.right-align` | `text-align:start/end` + `justify-content` | `text-left justify-start` / `text-right justify-end` |
| `.middle-align` / `.top-align` / `.bottom-align` | `display:flex` + `align-items` | `flex items-center` / `flex items-start` / `flex items-end` |
| `.circle` / `.round` | `border-radius:50%` / `2rem` | `rounded-full` / `rounded-card` |
| `.transparent` | fond, ombre, couleur neutralisés en `!important` | `bg-transparent shadow-none text-inherit` |
| `.border` | `.0625rem solid var(--outline)`, fond transparent | `border border-outline` |
| `.absolute` / `.fixed` / `.bottom` | position / `inset-block-end:0` | `absolute` / `fixed` / `bottom-0` |
| `.center` / `.middle` | `50%` + `translate(-50%)` | `left-1/2 -translate-x-1/2` / `top-1/2 -translate-y-1/2` |
| `.no-margin` / `.no-padding` | posent des variables internes beercss | `m-0` / `p-0` |
| `.responsive` sur `main` | `flex:1; padding:.5rem; max-inline-size:75rem; margin:0 auto` | `.slide-main` |
| `.responsive` sur `img` | `block-size:3rem; inline-size:3rem; object-fit:cover` | `size-12 object-cover` — rarement l'intention réelle |
| `.max` générique | `max-inline-size:100%` | `max-w-full` |
| `.max` dans un `nav` | `flex:1` | `flex-1` |
| `.large-text` | `font-size:1rem` | `text-base` |
| `.small` / `.medium` sur `hN` | échelle de police Material | `text-*` explicite |
| `.medium` / `.extra` génériques | 2.5 rem carré, ou 50 %, ou 4 rem, selon le parent | au cas par cas : il n'y a pas d'équivalent |
| `.fixed` sur le footer | `position:sticky; inset:0; z-index:11` | `.slide-footer` |
| `.center-left`, `.center-right`, `.padding` | **aucune règle** | supprimées |
| `.main`, `.xlarge-*`, `.xsmall-*`, `.full-width`, `.poppins-*` | définies dans `style.css` du talk | jetons `@theme` et utilitaires nommés |

### L'entrée la plus dangereuse n'est pas une classe

```css
*+:is(address,article,blockquote,code,.field,fieldset,form,.grid,
      h1,h2,h3,h4,h5,h6,nav,ol,p,pre,.row,section,aside,table,.tabs,ul)
{ margin-block-start: 1rem }
```

C'est le rythme vertical de tout le deck, et il ne s'écrit nulle part dans les slides — donc rien ne le signale. Le jour où beercss part, **tous les éléments de toutes les slides se collent**. Il devient une classe `.slide-prose` (`[&>*+*]:mt-4`) posée par les layouts. Premier point à vérifier au rendu.

### Les motifs nommés

Quatre `@utility`, pour les motifs qui reviennent et coûteraient une douzaine d'utilitaires chacun :

- `.stage` — la scène elle-même ;
- `.slide-main`, `.slide-header`, `.slide-footer`, `.slide-prose` — le châssis que les layouts posent ;
- `.contact-row` — la ligne « avatar rond + pseudo », répétée 12 fois dans `demoit.md:127-143` ;
- `.card` — le bloc `round extra` / `border` des slides de contenu.

## Changements côté Go

| Fichier | Changement |
|---|---|
| `deck/directive/transform.go:98-108` | émet `grid grid-cols-12 gap-4` et `col-span-{N}` au lieu de `grid` et `s{N}` ; `columnCount` reste 12 (largeur de la grille Tailwind) |
| `deck/directive/transform.go:28,33,86,92` | commentaires et message d'erreur qui nomment beercss, réécrits |
| `deck/directive/render.go:119-126` | `splitAttributes` émet `h-stage-{taille}` au lieu de `{height}-height` ; **l'API d'auteur ne change pas** (`split{height=xlarge}`, `height:` en frontmatter) |
| `deck/talk.go:11-18` | nouveau bloc `Theme struct{ Stage, Dark bool }` `\`yaml:"theme"\`` |
| `handlers/step.go:41-51,89` | champs `Page.Stage` et `Page.Dark`, remplis via `deck.LoadTalk` |
| `main.go:63` | route `/tailwind.css` → `handlers.Static` |
| `handlers/code.go` | **inchangé** : le `style=` arrive déjà de l'URL |

Les tests touchés : `deck/directive/transform_test.go:90` (et ses attentes de classes), `directive_test.go`, `talk_test.go` (bloc `theme:`), plus un nouveau cas pour `h-stage-*`.

### Le seul endroit où préserver l'existant est impossible

`render.go` garde l'API `height=`, mais les valeurs derrière changent. Aujourd'hui `.xlarge-height{block-size:64rem}` vaut 3686 px à 57.6 px/rem — 3,4 fois la hauteur de scène — et seule la requête `@media (max-height:1024px)` le ramène à 45 rem, soit 2592 px, encore le double. Ces valeurs sont incohérentes ; les jetons `--height-stage-{sm,md,lg,xl}` sont recalibrés sur la scène de 1080 px, à l'œil, sur les slides `split` du deck. C'est un changement de rendu assumé, pas une régression.

## Périmètre de fichiers

**Moteur** — `handlers/resources/index.tmpl.html`, `grid.tmpl.html`, `speakernotes.tmpl.html` ; `deck/layouts/` : `bare`, `content`, `cover`, `default`, `quote`, `split`, `partials` ; `deck/directive/transform.go`, `render.go` ; `deck/talk.go` ; `handlers/step.go` ; `main.go` ; nouveaux `styles/demoit.src.css`, `styles/tokens.css`, `handlers/resources/demoit.css`, `hack/css.sh` ; `.gitignore` (`.tools/`).

**Talk** — `impact-framework/demoit.md`, `demoit-en.html` ; `.demoit/layouts/cover.html`, `default-h3.html` ; `.demoit/style.css` ; `.demoit/js/demoit.js` ; `.demoit/talk.yml` ; nouveaux `.demoit/tailwind.src.css`, `.demoit/tokens.css`, `.demoit/tailwind.css`, `.demoit/fonts/` ; suppression de `.demoit/poppins.css` ; une ligne d'en-tête dans `demoit.html`, qui devient aussi une référence pré-Tailwind.

**Documentation** — `CLAUDE.md` décrit beercss dans cinq sections (grille 12 colonnes de `split{cols=}`, classes des slides, contenu de `.demoit/`, defaults des layouts, styling des deux talks). Elle deviendrait fausse.

## Le sort de `sample/`

L'opt-in protège `sample/` de la scène fixe et du commutateur. Il ne peut pas le protéger de tout, et il faut le dire franchement : **`index.tmpl.html` est embarqué et partagé, donc `sample/` perd la feuille de base de beercss le jour où elle part** — les tailles de `h1`…`h6`, le `font-size` et la `line-height` du `body`, le `display:flex` de `header`/`footer`, et la règle de rythme vertical `*+:is(…)`. À la place il reçoit le preflight de Tailwind et les jetons du moteur. Son rendu change donc, quoi qu'on fasse.

Ce n'est pas un choix évitable, c'est la conséquence directe de la décision 7. Deux conséquences pratiques :

1. `sample/` n'utilisant presque aucune classe beercss (`.mermaid` et `.logos`, toutes deux définies dans son propre `style.css`), le delta attendu est celui de la feuille de base, pas celui de la mise en page.
2. Si le preflight le casse visiblement, les quelques règles de base qui lui manquent sont ajoutées à `sample/.demoit/style.css` — un fichier qui appartient au talk, donc le bon endroit pour ça. Pas de modernisation, pas de scène, pas de thème : juste ce qu'il faut pour qu'il rende et navigue.

## Vérification

`go test ./...` couvre `deck`, `deck/directive`, `highlight`, `vscode` — il ne dit rien d'un rendu. La vérification est donc double.

**Automatique.** `go test ./...`, `gofmt -l .`, `golangci-lint run -c golangci.yml`. Nouveaux cas : les classes émises par `transform.go` et `render.go`, le bloc `theme:` de `talk.yml` et ses valeurs par défaut (absent = les deux booléens à faux).

**Visuelle, par comparaison.** Captures de référence des slides **sur le HEAD actuel, avant toute modification** — le worktree permet de servir l'avant et l'après simultanément. Puis comparaison slide par slide. C'est la méthode qui a servi à la migration Markdown de ce même deck.

Contrôles ciblés, une slide de chaque famille : `::term`, `::vscode`, `::browser`, `::code`, `split{cols=}` et `:::grid`/`:::col` écrits à la main, mermaid, `layout: quote`, `layout: cover`, `layout: bare`. Routes : `/`, `/N`, `/last`, `/grid`, `/pdf`, `/speakernotes`. Combinaisons : deux thèmes × trois profils. Et `sample/`, servi tel quel : il doit rendre, naviguer et faire fonctionner ses composants — au delta de feuille de base près, décrit plus haut.

## Hors périmètre

- Refonte esthétique des layouts (décision 8) : les jetons la préparent, elle ne se fait pas ici.
- `vs-code` et `web-browser` suivant le thème sombre.
- Vendorer mermaid, dernier point de dépendance réseau à l'exécution.
- `impact-framework/demoit.html` : conservé intact comme référence, une ligne d'en-tête ajoutée.
- `sample/` : vérifié, non modernisé — ni scène, ni thème, ni Tailwind dans ses slides.
