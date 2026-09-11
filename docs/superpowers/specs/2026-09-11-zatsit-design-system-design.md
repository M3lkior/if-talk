# Alignement du rendu des slides sur la nouvelle identité zatsit.fr

Date : 2026-09-11
Statut : design validé, prêt pour plan d'implémentation

## Problème

Le socle Tailwind posé le 2026-09-10 s'arrêtait volontairement au châssis : sa décision 8 disait « socle seulement, pas de refonte esthétique — les jetons de design posés ici rendent l'embellissement ultérieur trivial et réversible ». C'est cet embellissement.

Entre-temps zatsit.fr a été refait, et le deck ne lui ressemble plus. L'écart se mesure fichier contre fichier, `impact-framework/.demoit/tokens.css` contre `corporate/src/styles/global.css` du dépôt `zats-websites` :

**Le mode sombre ne dit pas la même chose.** Le site **inverse la teinte** : primaire bleu `#0f15fd` en clair, primaire **or `#f1be51`** en sombre, avec une famille d'accents ambre / rouge / orange et un dégradé de logo or → ambre → rouge. Le deck, lui, se contente d'éclaircir le bleu en `#7c81ff`. C'est l'écart le plus visible et le seul qui, à lui seul, fait dire « ce n'est pas le nouveau zatsit ».

**Les neutres viennent de deux familles différentes.** Le site est sur `slate` (`#1c1e21`, `#475569`, `#cbd5e1`, `#334155`), le deck sur une famille neutre chaude (`#3d4043`, `#6b7075`, `#d5d7da`, `#3a3d42`). Aucun des deux n'est faux, mais côte à côte ça se voit.

**Il manque tout le vocabulaire décoratif.** Le site a une grille de points, un mesh gradient, des blobs flottants, du verre dépoli, un dégradé de texte avec halo, des cartes dont la bordure passe à la primaire au survol. Le deck n'a rien de tout ça, et son `--radius-card` vaut `2rem` là où le site plafonne à `0.75rem` : les cartes du deck sont visiblement plus rondes que celles du site.

**Le logo n'est pas celui du site.** La cover empile deux images à `25rem` dont `zatsit_logo_noir.svg`, un aplat noir, sur une plaque en dégradé bleu → blanc (`.title` dans `style.css`) avec un titre forcé en `style="color: white"`. Le site, lui, inline `assets/icons/logo-full.svg` et le remplit d'un dégradé piloté par des jetons.

Objectif : que le deck lise comme le site, en clair et en sombre, sans rien casser du châssis ni des autres talks.

## Décisions

| # | Décision | Motif |
|---|---|---|
| 1 | **Inversion de teinte reprise telle quelle** : bleu en clair, or en sombre | C'est la signature du nouveau site. La demi-mesure (garder le bleu partout) donne un deck qui ne ressemble à rien de précis |
| 2 | Décor **statique** sur les slides de travail, **complet** sur `cover` et `quote` | Le site lui-même ne décore que son hero. Et ce deck fait tourner des tty et des iframes live pendant 45 min : un `blur(80px)` animé en permanence à côté est un coût permanent pour un effet qu'on ne regarde pas |
| 3 | Le décor est une propriété de **la scène**, pas des slides | Deux pseudo-éléments sur `.stage` : aucun markup de slide ne bouge, aucune layout n'a à le savoir, et le décor ne peut pas décaler une géométrie |
| 4 | `data-display="projector"` **éteint le décor** | Reste dans le contrat des profils d'affichage — de la lisibilité, jamais de la géométrie, donc ça ne peut rien déplacer |
| 5 | Le logo de la cover est **inliné**, pas en `<img>` | La def `#logo-full-gradient` et la référence `url(#…)` doivent être dans le même document. C'est exactement ce que fait `Hero.astro` via `<LogoFull />`. Un `<img>` ne peut pas être recoloré par le thème |
| 6 | Seuls **les mots accentués** du titre passent en dégradé | Le pattern du hero (« La tech *augmentée* au service de l'*impact* »), et le deck a déjà la convention `em` = accent. Le reste du titre garde son contraste plein, donc le profil `projector` a encore prise dessus |
| 7 | Titre et sous-titre de la cover remontent dans **`talk.yml`** | Demandé. La cover se compose alors seule, et `Talk.Title` était déjà parsé sans que rien ne le lise |
| 8 | Logos partenaires : **deux listes parallèles**, bascule en CSS | Le plus petit mécanisme qui marche. Un talk qui ne déclare pas `logosDark` ne change pas d'un pixel |
| 9 | **Toute la palette runtime vit dans le moteur** ; le `tokens.css` du talk la reflète pour la génération des utilitaires | Mesuré, pas supposé : le `tokens.css` d'un talk est importé en `theme(reference)` et **n'émet rien**. Voir « Qui émet réellement la palette » ci-dessous — c'est la contrainte qui commande ce découpage |

## Où va quoi

Le découpage en trois feuilles posé par la migration Tailwind ne bouge pas. Ce qui change, c'est leur contenu.

| Fichier | Ce qu'il gagne |
|---|---|
| `styles/tokens.css` | **la palette claire complète** (neutres, primaire, accents, stops du dégradé de logo), `--radius-card` à `0.75rem`, jetons de décor par défaut **inertes** |
| `styles/demoit.src.css` | **la palette sombre complète** dans `.dark`, machinerie du décor sur `.stage`, extinction du décor sous `projector` |
| `styles/components.css` | l'utilitaire `slide-hero` que `cover` et `quote` posent |
| `deck/layouts/cover.html`, `quote.html` | la classe `slide-hero` |
| `deck/layouts/partials.html` | la double émission des logos clair / sombre |
| `deck/talk.go` | `Subtitle`, `LogosDark`, et le rendu markdown de `Title` / `Subtitle` |
| `impact-framework/.demoit/tokens.css` | le **reflet** de la palette claire du moteur, pour que les utilitaires du talk se génèrent contre le bon vocabulaire |
| `impact-framework/.demoit/style.css` | `.logo-full-gradient-glow`, le traitement `em` sous `slide-hero`, Poppins 600 ; **perd** la plaque `.title` |
| `impact-framework/.demoit/layouts/cover.html` | le SVG inliné, la composition logo → titre → sous-titre |
| `impact-framework/.demoit/talk.yml` | `title`, `subtitle`, `logosDark` |
| `impact-framework/demoit.md` | la slide de cover perd son contenu, devenu inutile |

## 1. La palette

### Qui émet réellement la palette

Contrainte mesurée, contre-intuitive, et qui commande tout ce qui suit : **le `tokens.css` d'un talk n'émet aucune variable au navigateur.**

`tailwind.src.css` l'importe en `theme(reference)`, ce qui donne le vocabulaire au générateur d'utilitaires sans rien produire. Vérification : `impact-framework/.demoit/tailwind.css` ne contient **aucun** bloc `:root`, et la seule occurrence de `--color-main` y est un repli mort à l'intérieur d'un `var(…, #0f15fd)`. Toutes les valeurs qui s'appliquent vraiment viennent du `:root` et du `.dark` de `handlers/resources/demoit.css`.

Si la palette du talk fonctionne aujourd'hui, c'est donc uniquement parce que `impact-framework/.demoit/tokens.css` **duplique à l'identique** `styles/tokens.css`. Faire diverger les deux ne repeindrait rien.

Deux issues ont été essayées et écartées :

- **Faire émettre le talk** en passant l'import de son `tokens.css` en import simple. Ça marche — le `:root` apparaît bien dans `tailwind.css` — mais **ça casse le mode sombre**. Le bloc émis est `:root, :host`, hors couche de cascade, spécificité `(0,1,0)` ; le `.dark` du moteur est lui aussi hors couche et de spécificité `(0,1,0)`. À spécificité égale, l'ordre des feuilles tranche, et `/tailwind.css` est chargé **après** `/demoit.css` : le `:root` clair gagnerait sur le `.dark`. Mesuré sur un build réel, pas déduit.
- **Mettre la palette runtime dans le `style.css` du talk**, chargé en dernier. Correct en cascade, mais impose de maintenir chaque valeur à deux endroits — une dans `tokens.css` pour la génération, une dans `style.css` pour le rendu — avec zéro signal quand les deux divergent.

D'où la décision 9 : **la palette runtime reste là où elle est déjà, dans le moteur.** `styles/tokens.css` porte le clair, le `.dark` de `styles/demoit.src.css` porte le sombre, et `impact-framework/.demoit/tokens.css` en reste le reflet — exactement le rapport qu'ils entretiennent aujourd'hui. Aucun mécanisme nouveau, et le `--color-main: #0f15fd` déjà présent dans le moteur montre que celui-ci porte de toute façon la couleur de la maison comme défaut.

> À reporter dans `CLAUDE.md` à l'implémentation : il affirme que « palette and per-talk theming live in `<folder>/.demoit/tokens.css` », ce qui est vrai de la génération des utilitaires et faux du rendu.

### Clair — `styles/tokens.css`, reflété dans celui du talk

Les deux fichiers restent **`@theme` et rien d'autre** : `tailwind.src.css` importe le second en `theme(reference)`, et Tailwind rejette un fichier référencé qui porte une autre règle.

| Jeton | Valeur | Origine dans `global.css` |
|---|---|---|
| `--color-main` | `#0f15fd` | `--color-primary` |
| `--color-main-dark` | `#f1be51` | `--color-primary` du bloc sombre |
| `--color-fg` | `#1c1e21` | `--color-text` |
| `--color-fg-muted` | `#475569` | `--color-text-muted` |
| `--color-surface` | `#ffffff` | `--color-bg` |
| `--color-surface-raised` | `#f8fafc` | `--color-surface` / `--card-bg` |
| `--color-outline` | `#cbd5e1` | `--color-border` |
| `--color-accent-1` | `#06b6d4` | `--color-accent-1` |
| `--color-accent-2` | `#3b82f6` | `--color-accent-2` |
| `--color-accent-3` | `#6366f1` | `--color-accent-3` |
| `--logo-gradient-start` | `#0F15FD` | idem |
| `--logo-gradient-mid` | `#0a0ecc` | idem |
| `--logo-gradient-end` | `#020466` | idem |
| `--radius-card` | `0.75rem` | `--radius-lg` |
| `--color-text-gradient` | `linear-gradient(to bottom, #0F15FD 0%, #020466 100%)` | idem |
| `--decor-dots` | `0.08` | opacité de `.dot-grid` |
| `--decor-mesh` | `0.30` | opacité de `.hero-gradient` |

`--color-void` (le létterboxage autour de la scène) reste `#000000` : il n'a pas d'équivalent sur le site, et c'est voulu qu'il soit neutre.

`--color-text-gradient` ne peut pas être un jeton `@theme` exploitable tel quel : c'est un `linear-gradient`, pas une couleur, et il change aussi de **direction** entre les deux modes (`to bottom` en clair, `135deg` en sombre). Il est déclaré comme custom property ordinaire à côté de la palette, dans le même `:root` et le même `.dark`, et n'est lu que par le traitement `em` de `slide-hero`.

### Sombre — le `.dark` du moteur

Le moteur porte déjà le bloc `.dark` (`styles/demoit.src.css:158`). Il prend la palette sombre du site, accents compris :

```css
.dark {
    --color-fg: #e3e3e3;
    --color-fg-muted: #94a3b8;
    --color-surface: #1b1b1d;
    --color-surface-raised: rgba(255, 255, 255, 0.05);
    --color-outline: #334155;
    --color-main: var(--color-main-dark);

    --color-accent-1: #f59e0b;
    --color-accent-2: #ef4444;
    --color-accent-3: #f97316;

    --logo-gradient-start: #f1be51;
    --logo-gradient-mid:   #f59e0b;
    --logo-gradient-end:   #ef4444;

    /* et le mesh baisse, comme sur le site */
    --decor-mesh: 0.16;

    /* chrome des composants : --dm-* réaccordés sur ces neutres */
}
```

`--color-main` continue de passer par `--color-main-dark`, l'indirection qui existe déjà. Les autres jetons sont réassignés directement : leur inventer à chacun un compagnon `-dark` coûterait six niveaux d'indirection pour un gain nul, puisque le moteur porte de toute façon les deux palettes.

### Poppins 600

Le site charge 400 / 600 / 700 ; le talk embarque 100 / 400 / 500 / 700, donc **pas de 600**. On ajoute `poppins-600-latin.woff2` et `poppins-600-latin-ext.woff2` dans `.demoit/fonts/`, avec les deux `@font-face` correspondants et les mêmes `unicode-range` que les autres poids.

Un talk doit rendre **sans réseau** : les fichiers sont commités, comme les huit autres. S'ils s'avéraient impossibles à récupérer au moment du build, le repli est de mapper `font-semibold` sur le 500 déjà présent — visible mais pas bloquant. À trancher à l'implémentation, pas maintenant.

## 2. Le décor

### Mécanique

Deux pseudo-éléments sur `.stage`, dans `styles/demoit.src.css`, à côté de la définition de la scène :

- `::before` — la grille de points : `radial-gradient(var(--color-main) 1px, transparent 1px)`, `background-size: 32px 32px`, opacité `0.08` ;
- `::after` — le mesh : quatre `radial-gradient` elliptiques aux accents, masqués vers le bas par `linear-gradient(to bottom, black 60%, transparent 100%)`.

Les deux sont en `position: absolute; inset: 0; z-index: 0`, et le contenu de la scène passe au-dessus. `.stage` est déjà `overflow: hidden`, donc rien ne déborde.

### Opt-in, et extinction

Deux jetons pilotent l'ensemble, un par couche, parce que les deux ne varient pas ensemble — le site baisse son mesh en sombre et laisse sa grille de points inchangée :

```css
:root { --decor-dots: 0; --decor-mesh: 0; }   /* inerte par défaut */
```

Le moteur les active aux valeurs du site, `--decor-dots: 0.08` et `--decor-mesh: 0.30`, `.dark` ramenant le mesh à `0.16`. Les valeurs par défaut ci-dessus sont ce que voit un talk qui n'a pas la scène — et **`sample/`, qui n'a même pas de `talk.yml`, ne rend pas la scène du tout** : son `#app` porte `.no-stage`, donc les pseudo-éléments de `.stage` n'existent pas pour lui. Le décor lui est inaccessible par construction, pas seulement éteint.

Le profil projecteur éteint les deux :

```css
:root[data-display="projector"] { --decor-dots: 0; --decor-mesh: 0; }
```

C'est légitime au regard du contrat des profils — « un profil redéfinit des jetons de lisibilité et ceux-là seuls, jamais une valeur de géométrie » : l'opacité d'une couche purement décorative ne peut déplacer aucune boîte.

### Renforcement sur `cover` et `quote`

Ces deux layouts posent `slide-hero` sur leur `<main>`, et la scène réagit :

```css
.stage:has(main.slide-hero) { --decor-dots: 0.10; --decor-mesh: 0.45; }
.dark .stage:has(main.slide-hero) { --decor-mesh: 0.26; }
```

Les valeurs exactes sont à ajuster à l'œil sur la scène — c'est le seul réglage de cette spec qui se décide devant le rendu et pas sur le papier.

`:has()` plutôt qu'une classe sur `.stage` : le layout ne peut pas atteindre son propre conteneur, qui est écrit dans `index.tmpl.html`. Le sélecteur est supporté par tous les navigateurs visés, chromedp de `/pdf` compris.

**Attention au piège des classes par défaut.** Un layout possède ses classes de `<main>` par défaut, et une slide qui déclare `class:` **remplace la liste en bloc**. La slide de cover d'`impact-framework` déclare aujourd'hui `class: slide-main flex flex-col items-center justify-center text-center title` — donc elle écraserait `slide-hero`. Cette ligne disparaît de toute façon (la plaque `.title` n'existe plus), mais c'est le genre d'oubli qui ne se voit qu'à l'écran.

### Ce que le décor ne touche pas

`.mermaid` reste exactement comme il est. Le CLAUDE.md est formel : mermaid mesure son conteneur pour dimensionner ses colonnes, toute marge intérieure ou tout `display` sur le `svg` rogne les libellés — un défaut connu et non résolu sur cette figure. Le décor est derrière, jamais dedans.

## 3. La cover

### L'asset

`corporate/src/assets/icons/logo-full.svg` est copié en `impact-framework/.demoit/images/logo-full.svg`. Mêmes 7 `path` et même `viewBox="0 0 523.94 426.01"` que le `zatsit.svg` déjà présent, mais avec `fill="currentColor"` et la def de dégradé aux stops pilotés par `--logo-gradient-*`.

### Le layout

`impact-framework/.demoit/layouts/cover.html` inline le SVG et compose :

```
<svg class="logo-full-gradient-glow" …>  ← les 7 path + la def #logo-full-gradient
{{ .Talk.Title }}                         ← rendu markdown, les *mots* deviennent <em>
{{ .Talk.Subtitle }}                      ← rendu markdown, en --color-fg-muted
{{ logos clair / sombre }}
```

et `style.css` reprend les deux règles du site mot pour mot :

```css
.logo-full-gradient-glow path { fill: url(#logo-full-gradient); }
.logo-full-gradient-glow {
    filter: drop-shadow(0 0 40px color-mix(in srgb, var(--color-main) 40%, transparent));
}
```

Le sélecteur vise `path` et non l'élément : c'est ce que fait le site, et c'est correct dès lors que le SVG est un vrai élément du document. **Un `<use>` vers un `<symbol>` ne marcherait pas** — un sélecteur du document n'atteint pas l'arbre fantôme d'un `<use>`, les `path` retomberaient sur le noir par défaut. C'est l'erreur que les maquettes de brainstorming ont faite, corrigée là-bas en posant le `fill` sur l'élément `<svg>` lui-même, qui lui est hérité.

### `em` : deux traitements

Le deck a déjà `em { font-style: normal; color: var(--color-main) }`. On garde ça partout, **sauf** sous `slide-hero`, où `em` prend le dégradé et le halo :

```css
.slide-hero em {
    background: var(--color-text-gradient);
    background-clip: text;
    -webkit-text-fill-color: transparent;
    filter: drop-shadow(0 0 30px color-mix(in srgb, var(--color-main) 50%, transparent));
}
```

Un `em` sur une slide de travail reste donc un aplat d'accent lisible ; sur la cover et les citations, c'est la signature du site.

### Ce qui disparaît

- la plaque `.title` de `style.css` (dégradé bleu → blanc, devenue sans objet) ;
- dans `impact-framework/demoit.md`, le `<h1 class="m-0" style="color: white">` et le `<h3 style="color: white">` de la slide de cover, avec sa ligne `class:`. Ce sont les **deux seules couleurs en dur du deck** — vérifié par `grep`, le reste des occurrences est dans le `classDef` mermaid, qui garde les siennes.

## 4. Titre et sous-titre dans `talk.yml`

`deck/talk.go` gagne un champ, et deux champs changent de type :

```go
type Talk struct {
    Title    template.HTML  // rendu markdown au chargement
    Subtitle template.HTML  // idem, nouveau
    …
}
```

Le rendu passe par `renderTitle` (`deck/deck.go:175`), la fonction qui convertit déjà le `title:` d'une slide par le même pipeline markdown et le réduit à un seul paragraphe. C'est ce qui fait que `*Impact Framework*` devient un `<em>`.

`Title` est aujourd'hui un `string` lu par `talk_test.go:39` ; ce test devra comparer un `template.HTML`. Aucun layout embarqué ni de talk ne lit `Talk.Title` à ce jour, donc le changement de type ne casse aucun rendu.

Dans `talk.yml` :

```yaml
title: A la découverte d'*Impact Framework*.
subtitle: 10 avril 2025 — Ludovic Dussart
```

## 5. Logos partenaires, clair et sombre

`logo.jpg` est un JPEG à fond blanc opaque : sur une slide sombre c'est un rectangle blanc. On ne détoure rien — on prévoit la mécanique, les fichiers viendront.

`deck/talk.go` :

```go
LogosDark []string `yaml:"logosDark"`
```

`deck/layouts/partials.html` (le header) et les deux `cover.html` émettent les deux jeux :

```
{{ range .Talk.Logos }}<img class="… {{ if $.Talk.LogosDark }}dark:hidden{{ end }}" src="{{ . }}">{{ end }}
{{ range .Talk.LogosDark }}<img class="… hidden dark:block" src="{{ . }}">{{ end }}
```

Le `{{ if }}` est le cœur du mécanisme : **sans `logosDark`, aucune classe conditionnelle n'est émise et le rendu est identique à aujourd'hui, dans les deux thèmes.** C'est ce qui protège `sample/` et tout talk existant.

`hidden`, `dark:hidden` et `dark:block` sont des utilitaires Tailwind standards, vus par le scanner du moteur dans `deck/layouts/` et par celui du talk via son `@source "../"`. Aucun `@source inline` supplémentaire n'est requis pour eux — contrairement à `slide-hero`, qui est un `@utility` et doit donc rejoindre la ligne `@source inline(…)` du moteur, faute de quoi il ne serait émis que là où un fichier scanné l'utilise.

## Vérification

La suite Go et le pipeline CSS sont verts au départ : `go build ./...`, `go test ./...` (deck, deck/directive, handlers, highlight, vscode) et `hack/css.sh engine` + `hack/css.sh impact-framework` reproduisant les artefacts commités à l'octet près.

Le cache de build Go doit pointer ailleurs que `~/Library/Caches` — `GOCACHE="$TMPDIR/go-build"` — sinon le bac à sable refuse l'écriture.

Ce que la suite Go **ne dit pas** : ce que le navigateur peint. Les deux stylesheets sont des artefacts commités qu'aucun `go build` ne régénère, donc `hack/css.sh` fait partie de la vérification au même titre que les tests. Et chaque changement de rendu se regarde slide par slide sous `demoit --dev impact-framework`, dans les deux thèmes et sur les trois profils d'affichage — c'est la méthode par laquelle la migration markdown de ce talk avait elle-même été contrôlée.

Points à regarder en priorité, parce qu'ils sont invisibles pour les tests :

- la cover dans les deux thèmes — c'est la slide la plus retravaillée ;
- la figure mermaid en sombre, dont les libellés sont déjà rognés à droite : vérifier que le décor ne l'aggrave pas ;
- une slide `split` avec un tty et une iframe, décor actif, pour confirmer que rien ne passe au-dessus ni ne clignote ;
- `/pdf`, qui force `?theme=light` ;
- `/grid` et `/speakernotes`, qui suivent le thème par `BroadcastChannel`.

## Hors périmètre

- **Les blobs animés et le verre dépoli.** Écartés à la décision 2, pas oubliés. Si le décor statique déçoit à l'écran, ils se rajoutent ensuite sur `slide-hero` uniquement, sans toucher au reste.
- **`sample/`.** Il ne déclare pas de `talk.yml`, donc ni scène, ni thème sombre, ni décor. Il ne doit rien voir de ce chantier, et c'est un critère de vérification, pas une supposition.
- **`impact-framework/demoit-en.html` et `demoit.html`.** Le premier est un deck HTML localisé qui suit le chemin d'origine ; le second n'est jamais servi (`deck.resolve` préfère `demoit.md`) et est conservé comme référence de pré-migration. Aucun des deux n'est repeint ici.
- **Le défaut de rognage de mermaid.** Connu, documenté, non résolu, et hors sujet ici.
- **Le chrome des runtimes embarqués.** Le tty reste sombre en permanence (gotty fige son thème au démarrage), `vs-code` et `web-browser` ne suivent pas le thème. Inchangé.
