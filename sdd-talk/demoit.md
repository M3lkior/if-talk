---
layout: cover
class: responsive max center-align title
speakernotes: |
  S1 - Titre + toi. Qui je suis, pourquoi je parle de ça, la promesse de la
  demi-journée.
  Annoncer le format : 45 minutes d'introduction, puis atelier.
  Thèse unique à poser dès maintenant : l'IA a rendu le code abondant, le
  goulot s'est déplacé en amont (l'intention) et en aval (la vérification).
---

<h1 class="medium no-padding center-align top-align" style="color: white">Le développeur Product Engineer</h1>
<h3 style="color: white">Le développement piloté par les spécifications</h3>
<h5 style="color: white">Ludovic Dussart &mdash; Zatsit</h5>

---
layout: default
title: "Product engineer : le métier que le marché a inventé deux fois"
speakernotes: |
  S2 - Deux colonnes.
  Gauche, 2018-2024 : né sans IA, pour supprimer le coût des handoffs.
  PostHog (187 personnes, 47 équipes, 4,2 personnes en moyenne), incident.io,
  Jean-Michel Lemieux (ex-VP Eng Shopify), le Product Engineer Manifesto de
  Viljami Kuosmanen (~2024).
  Droite, 2025-2026 : généralisé par l'IA.
  Punchline : ce rôle n'a pas été créé par l'IA, il a été révélé par elle.
  Accroche possible : sur les job boards FR, "product engineer" désigne aussi
  un ingénieur produit R&D classique. Qui a cherché le terme ne trouve pas ce
  métier.
  Sources : PostHog Product Engineer Handbook (posthog.com/product-engineer) ;
  Product Engineer Manifesto (productengineer.org) ; Le Talent Club.
---

*[S2 &mdash; à produire]*

---
layout: default
title: "Pourquoi maintenant : le mécanisme"
speakernotes: |
  S3 - La spécialisation était une réponse à la complexité. Quand l'IA absorbe
  une partie de la complexité d'exécution, la division du travail perd sa
  justification : l'IA élimine la distance entre les disciplines.
  Illustrer avec le "one-person band" du Chief AI Officer de Pendo : au lieu de
  répartir entre engineering, produit et design, une seule personne livre,
  itère, collecte le feedback et itère de nouveau.
  Puis l'argument économique, le plus solide : la ressource rare n'est plus
  l'exécution, c'est le jugement.
  Chiffre optionnel, à manier avec précaution (petits volumes, Angleterre
  seulement) : 72 offres "Product Engineer" sur 6 mois au 5 janvier 2026,
  contre 20 un an avant. Directionnel, pas probant - le dire.
  Sources : Turing College, The Rise of the AI Product Engineer ; SFEIR, page
  concept Product Engineer ; itjobswatch.co.uk (Product Engineer, England).
---

*[S3 &mdash; à produire]*

---
layout: default
title: Le paradoxe de productivité
speakernotes: |
  S4 - Le personnage est posé, on montre son problème.
  Télémétrie sur plus de 10 000 développeurs : +21 % de tâches complétées,
  +98 % de PR mergées à l'échelle individuelle, métriques de delivery
  organisationnelles plates.
  LinearB, 8,1 millions de PR : le code généré par IA attend 4,6 fois plus
  longtemps sa première revue.
  METR : oscillation de 37 points sur 12 mois, d'un ralentissement de 19 % à
  une accélération de 18 %. La courbe en J.
  Conclusion DORA : l'IA ne répare pas une équipe, elle amplifie ce qui est
  déjà là.
  Sources : DORA 2026, The ROI of AI-assisted Software Development (Google) ;
  LinearB ; METR ; arXiv 2609.00252 pour le cadrage académique du paradoxe.
---

*[S4 &mdash; à produire]*

---
layout: default
title: Le vrai goulot n'est pas le code
class: main responsive large-height
speakernotes: |
  S5 - LA charnière du talk, avec S30. Trois temps, dans cet ordre.
  1. Le diagnostic : le temps de développement est rarement le vrai frein. Les
  goulots sont en amont, dans des exigences floues et des décisions produit
  faibles. Une IA qui rend le code moins cher ne répare rien de ça, elle
  l'expose : la partie du process qui absorbait le mou avance désormais plus
  vite que tout ce qui l'entoure.
  2. La phrase charnière. La lire à voix haute, lentement, et marquer un temps.
  C'est le seul moment du talk où le lien product engineer / SDD se formule en
  une phrase. Elle revient à l'identique en S30.
  3. Le plan en 4 temps, à l'oral uniquement, rien à l'écran : les prérequis,
  le SDD et son intérêt, le panorama des frameworks, découpler et accélérer.
  20 secondes, pas plus, puis enchaîner.
  Corollaire à dire, pas à écrire : l'ingénieur le plus utile n'est pas celui
  qui produit le plus de code, c'est celui qui opère bien sur tout le chemin
  d'un problème métier à un changement livré et fonctionnel.
  Source du diagnostic : hitechnology.io, The Rise of the Product Engineer
  (étude 100+ leaders) ; DORA 2026 (Google).
---

<div class="large-space"></div>

### Le temps de développement est rarement le vrai frein.

##### Les goulots sont en amont : des exigences floues, des décisions produit faibles.

<div class="small-space"></div>

#### Une IA qui rend le code moins cher ne répare rien de ça. Elle l'*expose*.

<div class="large-space"></div>

<blockquote class="charniere">
  <h4>Sans SDD, un product engineer n'est qu'un dev surchargé à qui on a retiré son PM.</h4>
  <h5>Il a besoin d'un artefact pour porter l'intention jusqu'à la machine. Cet artefact, c'est la <em>spec</em>.</h5>
</blockquote>

---
layout: default
title: L'échelle de délégation
speakernotes: |
  S6 - Autocomplétion, chat, agent, agent autonome.
  Le secteur passe de pratiques assistées comme le vibe coding, où l'assistant
  accélère un développeur isolé, vers l'Agentic Software Engineering où des
  agents autonomes reçoivent des tâches au niveau de l'objectif.
  Faire lever la main : qui est où aujourd'hui ? Ça calibre la salle pour toute
  la suite. Si la salle est majoritairement en bas de l'échelle, prévoir de
  transférer 3 minutes de la partie 3 vers les parties 1 et 2.
  Source : arXiv 2609.00252, Spec-Driven Development for Agentic Software
  Engineering.
---

*[S6 &mdash; à produire]*

---
layout: default
title: Comment marche vraiment un agent
speakernotes: |
  S7 - Le slide le plus utile du talk. Si celui-là passe, le reste coule.
  Pas de mémoire persistante, seulement une fenêtre de contexte.
  Tout le SDD découle de là : puisque les LLM n'ont pas de mémoire persistante,
  il faut créer des artefacts durables et versionnés (constitution.md, spec.md,
  plan.md) pour externaliser l'état du projet.
  Introduire le terme context engineering ici.
  Prendre le temps. C'est le slide à ne pas accélérer même en retard.
---

*[S7 &mdash; à produire]*

---
layout: default
title: Les modes de défaillance du vibe coding
speakernotes: |
  S8 - Concret, ils l'ont tous vécu. Raconter, ne pas lister.
  L'agent produit mille lignes qui ont l'air correctes et qu'on n'a pas
  demandées. On voulait un petit fix, on récupère un refactor, une logique
  "améliorée", trois nouveaux fichiers, et des tests qui passent parce qu'ils
  ne testent rien.
  Trois causes racines : l'intention ne vit que dans l'historique de chat ; la
  dérive ; l'output non vérifiable faute de critères d'acceptation.
  Ces trois causes annoncent les trois réponses du SDD - le dire.
---

*[S8 &mdash; à produire]*

---
layout: default
title: Rien de tout ça n'est nouveau
speakernotes: |
  S9 - Slide anti-objection. Frise : user story + critères d'acceptation,
  Gherkin/BDD, TDD, contract-first OpenAPI, ADR, DDD.
  Message : on a toujours écrit des specs. Ce qui change, c'est qu'elles
  doivent être lisibles par une machine et versionnées.
  Glisser le glossaire en encart sur ce slide plutôt qu'en slide séparé :
  spec, plan, task, constitution, gate, subagent, drift. Gain d'une minute.
---

*[S9 &mdash; à produire]*

---
layout: default
title: "DDD : pourquoi il compte plus qu'avant"
speakernotes: |
  S10 - Langage ubiquitaire, bounded contexts, context map, agrégats,
  invariants.
  L'argument : un nommage précis et cohérent affûte les prompts et empêche les
  agents de confondre des concepts distincts. Les agents ont accès à tout le
  dépôt et cherchent les noms qu'on leur donne ; ils peuvent naïvement essayer
  d'unifier des entités similaires, donc il faut expliciter qu'elles sont
  séparées pour une raison.
  Trois usages directs : langage ubiquitaire = désambiguïsation de l'agent ;
  bounded contexts = périmètre de contexte et de parallélisation ; invariants =
  critères vérifiables.
  Mise en garde qui compte : la valeur du modèle de domaine réside dans la
  compréhension partagée qu'il construit dans l'équipe, pas dans l'artefact.
  Une spec générée par IA que personne ne lit ne résout rien.
  Sources : threedots.tech, Domain-Driven Design matters more when AI writes
  your code ; arXiv 2605.01160 (section DDD/DDT par tiers de délégation) ;
  talk gitnation From Prompt Spaghetti to Bounded Contexts.
---

*[S10 &mdash; à produire]*

---
layout: default
title: "Spec-Driven Development : la définition"
speakernotes: |
  S11 - Le SDD traite les spécifications comme des contrats exécutables
  desquels les agents dérivent le code, en empêchant la dérive architecturale
  par une application automatisée plutôt que par de la documentation passive.
  La formule à afficher : the spec is the prompt.
  Source : Birgitta Böckeler, Understanding Spec-Driven Development
  (martinfowler.com) ; podcast Thoughtworks What is spec-driven development.
---

*[S11 &mdash; à produire]*

---
layout: content
title: Les 3 niveaux d'ambition
speakernotes: |
  S12 - Le cadre analytique le plus cité du domaine (Böckeler), à ne pas
  sauter.
  Spec-first : la spec guide l'implémentation puis est abandonnée, le code
  reste l'artefact maintenu. Point d'entrée pragmatique, là où sont la plupart
  des équipes qui démarrent.
  Spec-anchored : la spec persiste comme contrat vivant, versionnée et mise à
  jour. La cible réaliste 2026.
  Spec-as-source : la spec EST la source, le code est généré, jetable,
  régénérable. Vision long terme, encore largement expérimentale.
  Dire explicitement à quel niveau on propose de jouer : ça évite 80 % des
  malentendus dans la salle. C'est aussi la moitié de la réponse à la question
  "c'est le cycle en V déguisé ?".
  Source : Böckeler, martinfowler.com/articles/exploring-gen-ai/sdd-3-tools.html
---

*[S12 &mdash; à produire]*

---
layout: default
title: La boucle canonique
speakernotes: |
  S13 - Le diagramme à produire soi-même, et à réutiliser 3 fois dans le deck
  comme repère de progression.
  intention, clarification, spec (quoi/pourquoi + critères), plan (comment),
  tasks atomiques, implémentation, vérification, réconciliation spec/code.
  Les gates humains en rouge.
---

*[S13 &mdash; à produire]*

---
layout: default
title: Anatomie d'une bonne spec
speakernotes: |
  S14 - Le slide dont on se souviendra, à condition de le nourrir avec une
  spec réelle de mon propre code, en deux versions : floue vs complète, avec le
  résultat produit par l'agent dans les deux cas. C'est du travail à produire.
  Une bonne spec fixe les résultats attendus, les limites de périmètre, les
  contraintes, les décisions antérieures, la découpe en tâches et les critères
  de vérification.
  Punchline : tout ce que tu n'écris pas, l'agent l'invente.
  Source : github/spec-kit (lire les templates, pas le README).
---

*[S14 &mdash; à produire]*

---
layout: default
title: "Pourquoi ça marche, mécaniquement"
speakernotes: |
  S15 - Cinq mécanismes, pas cinq bénéfices marketing.
  1. Externalise l'état que le LLM n'a pas.
  2. Déplace la revue en amont : 200 lignes de markdown au lieu de 2000 lignes
  de diff.
  3. Rend l'output vérifiable.
  4. Rend le travail reprenable entre sessions et transférable entre agents.
  5. Permet le parallélisme.
  Appui : le SDD attrape les violations architecturales et la dérive de contrat
  d'API que les tests unitaires ne peuvent structurellement pas détecter, et il
  scale sur des agents parallèles en séparant le rôle qui implémente de celui
  qui vérifie.
---

*[S15 &mdash; à produire]*

---
layout: content
title: "SDD, TDD, BDD, cycle en V"
speakernotes: |
  S16 - Tableau compact, puis traiter frontalement la critique.
  Fowler soutient qu'une spécification utile dépend de l'apprentissage acquis
  pendant le développement, et que la clé de l'usage plein de l'IA est
  d'accélérer les boucles de feedback.
  La nuance à porter : Beck s'oppose à l'écriture de la spécification complète
  avant l'implémentation. Le niveau 2 traite la spec comme un document vivant
  tout au long de l'implémentation. La critique vise surtout le niveau 3 et
  tout workflow qui gèle les hypothèses.
  Sources : Kent Beck et Martin Fowler ; François Zaninotto, Waterfall Strikes
  Back (la critique à connaître pour tenir les questions).
---

*[S16 &mdash; à produire]*

---
layout: default
title: "Le coût, honnêtement"
speakernotes: |
  S17 - Slide de crédibilité. Trois coûts : cérémonie mal dimensionnée,
  tokens, dérive spec/code.
  L'anecdote parfaite : Böckeler a lancé Kiro sur une petite correction de bug
  et a obtenu quatre user stories et seize critères d'acceptation. Le SDD passe
  mal à l'échelle vers le bas : pour un null check d'une ligne, on n'a pas à
  toucher à la constitution.
  Sur la dérive : aucun de ces outils ne réconcilie automatiquement. Il faut le
  déclencher et relire la sortie, et c'est cette charge que la plupart des
  équipes esquivent jusqu'à ce que leurs specs aient six mois de retard.
  Ce slide porte la réponse à deux des cinq questions attendues : le coût en
  tokens, et qui maintient les specs quand elles divergent. Assumer que c'est
  le point faible non résolu du domaine.
  Source : Böckeler, martinfowler.com.
---

*[S17 &mdash; à produire]*

---
layout: default
title: "La carte par couches, et le niveau 0 : le mode plan"
speakernotes: |
  S18 - Deux choses sur un slide.
  D'abord pourquoi c'est confus : les outils opèrent à des couches différentes
  (définition des artefacts d'exigence, conversion en graphe de tâches,
  exécution du code, intégration IDE). Une cartographie communautaire
  recensait plus de 30 outils début 2026.
  Ensuite la baseline : le mode plan (Claude Code, Cursor). Plan éphémère,
  aucun artefact persisté, aucun gate, zéro installation. Excellent pour une
  tâche de 30 minutes, insuffisant dès qu'il y a plusieurs sessions, plusieurs
  agents ou une revue par un tiers.
  C'est le témoin de comparaison sur les 4 slides suivants - le dire
  explicitement, sinon le panorama devient un catalogue.
---

*[S18 &mdash; à produire]*

---
layout: default
title: GitHub Spec Kit
speakernotes: |
  S19 - CLI Python, environ 129 000 étoiles et 38 intégrations d'agents.
  Quatre commandes portent le workflow : /speckit.specify capture le contexte
  métier et les critères de succès, /speckit.plan traduit en décisions
  d'architecture, /speckit.tasks décompose en unités testables.
  Différenciateurs : une constitution définie une fois pour le projet dont
  chaque spec hérite, et des templates qui marquent les inconnues en
  NEEDS CLARIFICATION plutôt que de deviner.
  Faiblesses : pas d'étape de revue de code intégrée, et changer de direction
  implique de relancer les commandes concernées, chacune régénérant tout son
  document.
  Profil : greenfield, gates forts, verbeux.
  Sources : github/spec-kit ; ranthebuilder.cloud, I Tested Three Spec-Driven
  AI Tools. Chiffre d'étoiles à revérifier avant de figer.
---

*[S19 &mdash; à produire]*

---
layout: default
title: "OpenSpec : l'anti-cérémonie"
speakernotes: |
  S20 - L'empreinte la plus légère : on écrit des delta specs, uniquement ce
  qui change.
  Installation en 5 minutes contre 30, pas de Python, sortie d'environ 250
  lignes contre 800, conçu pour les bases de code existantes, 3 commandes IA
  contre 8, mise à jour de n'importe quel artefact à tout moment sans phase
  gates rigides.
  Workflow OPSX : explore, propose, apply, sync, archive.
  Contrepartie : pas de gates de revue entre phases.
  C'est aussi la réponse à la question "et le legacy ?".
  Source : github.com/Fission-AI/OpenSpec (docs/opsx.md et docs/commands.md) ;
  reenbit.com/bmad-vs-spec-kit-vs-openspec pour les chiffres de comparaison.
---

*[S20 &mdash; à produire]*

---
layout: default
title: "Superpowers : la méthode comme artefact"
speakernotes: |
  S21 - Le cas le plus intéressant pédagogiquement. Ce n'est pas un dépôt de
  specs, c'est une méthodologie livrée comme skills.
  Plugin qui impose un workflow structuré avant qu'une seule ligne de code ne
  soit écrite : brainstormer d'abord, isoler sa branche, écrire un plan
  détaillé, exécuter. Chaque étape conditionne la suivante.
  Le plan découpe le travail en tâches de 2 à 5 minutes avec chemins de
  fichiers exacts, commandes exactes et code complet. Des subagents
  implémentent chaque tâche, avec une revue en deux étapes.
  Pratiques imposées : cycles TDD red-green-refactor où les tests doivent
  échouer avant l'implémentation, méthodologie de debug en quatre phases
  exigeant l'investigation de la cause racine avant tout correctif, sessions de
  brainstorming socratique.
  L'angle à souligner devant des devs : c'est une réponse directe au problème
  de discipline. Quelqu'un écrit une bonne skill specify-plan-implement,
  l'utilise une semaine, puis retourne discrètement au prompting non structuré
  dès qu'une deadline approche.
  Anecdote : Jesse Vincent, créateur de RT, contributeur de Perl 5, auteur du
  client mail K-9, a quasiment cessé de coder lui-même.
  Source : claude.com/plugins/superpowers ; obra/superpowers-marketplace.
---

*[S21 &mdash; à produire]*

---
layout: default
title: "BMAD-METHOD : le poids lourd"
speakernotes: |
  S22 - Simule une équipe agile complète via des personas nommés : Analyst,
  Product Manager, Architect, UX Designer, Scrum Master, Dev, et BMad Master
  comme orchestrateur.
  L'agent Scrum Master crée des fichiers de story détaillés portant le contexte
  architectural, les guidelines d'implémentation et les critères de test pour
  l'agent Dev.
  En v6 : intelligence adaptative à l'échelle qui ajuste la profondeur de
  planification du bugfix au système d'entreprise, 19 agents spécialisés, trois
  pistes - Quick Flow avec tech spec seule, BMad Method avec PRD + architecture
  + UX, Enterprise pour la conformité.
  Source : github.com/bmad-code-org/BMAD-METHOD (v6). Version à revérifier.
---

*[S22 &mdash; à produire]*

---
layout: content
title: Le tableau de décision
speakernotes: |
  S23 - Le chiffre qui fait rire la salle : sur un même build de dashboard CRM,
  la même tâche a pris 12 minutes avec OpenSpec, 90 minutes avec Spec Kit et
  5 h 30 avec BMAD.
  Puis la grille : pour la plupart des équipes travaillant sur du code
  existant, OpenSpec offre le meilleur équilibre vitesse/flexibilité ; pour les
  nouveaux projets avec des rôles clairs, Spec Kit apporte structure et
  documentation ; pour la complexité d'échelle entreprise, BMAD gère
  l'orchestration multi-agents.
  L'honnêteté qui fait la différence : BMAD est excellent, mais aussi coûteux
  et disproportionné pour la plupart du travail hebdomadaire d'ingénierie.
  Mentionner en une ligne, sans slide : Kiro, Tessl, cc-sdd, Antigravity, et la
  convergence AGENTS.md / Agent Skills comme lingua franca émergente.
  Source des durées : reenbit.com/bmad-vs-spec-kit-vs-openspec.
---

*[S23 &mdash; à produire]*

---
layout: default
title: Ce sur quoi ils sont tous d'accord
speakernotes: |
  S24 - Slide de sortie de partie. Quatre primitives universelles :
  règles/constitution, spec, plan, tasks, plus des gates humains.
  Malgré des approches différentes, tous ces frameworks s'accordent sur un
  point : l'humain reste dans la boucle, mais pas pour tout.
  Conclusion libératrice, c'est le message à laisser : tu peux commencer demain
  avec trois fichiers markdown et un gate. Le framework est une optimisation,
  pas un prérequis.
---

*[S24 &mdash; à produire]*

---
layout: content
title: Les 5 découplages
speakernotes: |
  S25 - 1. Le quoi / le comment : la spec devient l'interface entre jugement
  humain et exécution machine.
  2. La conception / l'exécution dans le temps : tu spécifies maintenant, les
  agents exécutent pendant que tu fais autre chose. Ton temps de clavier n'est
  plus le facteur limitant.
  3. La revue / le code : revoir l'intention (200 lignes lisibles) au lieu du
  diff (2000 lignes qu'on n'a pas écrites).
  4. Les développeurs entre eux : specs comme contrats + bounded contexts =
  agents parallèles sur des tranches indépendantes, sans collision (worktrees,
  subagents).
  5. L'humain / la session : spec et plan comme mémoire durable. Reprise après
  crash de session, transfert entre agents, onboarding.
  Annoncer que le n°3 est celui sur lequel on s'arrête.
---

*[S25 &mdash; à produire]*

---
layout: default
title: "Le découplage qui compte : la revue"
speakernotes: |
  S26 - Pourquoi le découplage n°3 est LE gain.
  Le rapport DORA 2026 montre que l'IA n'élimine pas les goulots, elle déplace
  souvent le problème à l'étape suivante : si l'équipe augmente sa production
  de code mais continue à revoir les changements de la même façon, avec un
  contexte limité et peu de signaux de risque clairs, une partie de la vitesse
  gagnée en développement se perd en validation.
  Le SDD attaque ce point précis en donnant au reviewer l'intention explicite
  contre laquelle juger.
  Faire le lien avec le 4,6x de S4.
  Source : DORA 2026, The ROI of AI-assisted Software Development (Google).
---

*[S26 &mdash; à produire]*

---
layout: default
title: "Les chiffres d'accélération, et leurs limites"
speakernotes: |
  S27 - Le meilleur point de données à l'échelle d'une équipe : Mercari, place
  de marché japonaise d'environ 22 millions d'utilisateurs mensuels, a publié
  sa méthodologie interne "Agent Spec-Driven Development" en décembre 2025.
  Après six mois, elle rapportait un gain de vitesse de 150 % sur sa baseline
  traditionnelle et de 80 % sur le prompting IA en format libre.
  Compléments : GitHub rapporte que les équipes utilisant Spec Kit livrent avec
  environ un ordre de grandeur moins de cycles de régénération from scratch ;
  AWS documente des cas où des fonctionnalités de 40 heures ont été livrées en
  moins de 8 heures de temps humain quand elles étaient d'abord rédigées comme
  specs.
  Dire les limites à voix haute : ce sont des chiffres auto-rapportés, sur un
  seul contexte d'équipe.
  Sources : github.com/ianhxu/agentic-engineering-field-study
  (04-spec-driven-development.md, cas Mercari) ; GitHub ; AWS.
---

*[S27 &mdash; à produire]*

---
layout: default
title: Où le SDD n'accélère pas
speakernotes: |
  S28 - Indispensable pour la crédibilité, et c'est le slide qui fait gagner la
  salle.
  DORA 2026 montre des gains de 35 à 40 % sur les tâches simples mais moins de
  10 % sur du code legacy complexe.
  Donc : petits changements, spikes exploratoires, domaine encore flou, legacy
  sans tests.
  Message : le SDD est un investissement dont le retour dépend de la taille du
  changement.
  Source : DORA 2026, The ROI of AI-assisted Software Development (Google).
---

*[S28 &mdash; à produire]*

---
layout: default
title: Quoi mesurer
speakernotes: |
  S29 - Pas les lignes de code.
  Lead time jusqu'à la prod, temps de revue, taux de rework, change failure
  rate, et surtout la part du rework due à une intention mal comprise : c'est
  la métrique que le SDD prétend améliorer.
  Chiffre à citer : seules 7,3 % des équipes ont un taux de rework sous 2 %, ce
  qui révèle une taxe cachée sur la productivité que la génération de code par
  IA peut facilement aggraver.
  Source : DORA 2026, The ROI of AI-assisted Software Development (Google).
---

*[S29 &mdash; à produire]*

---
layout: default
title: "Boucle fermée : l'artefact du product engineer"
class: main responsive large-height
speakernotes: |
  S30 - La résolution de S5. Même layout, même phrase, même mise en forme :
  la salle doit voir la boucle se fermer sans qu'on l'explique.
  1. Le product engineer avait besoin d'un artefact pour porter l'intention. On
  vient de passer 35 minutes à décrire cet artefact et son outillage.
  2. Relire la phrase charnière. C'est le même bloc qu'en S5, au mot près.
  Marquer le même temps d'arrêt qu'au début.
  3. La nouvelle pile de compétences : cadrage de problème, modélisation de
  domaine, écriture de critères vérifiables, arbitrage, design de la
  vérification. Et la vitesse de frappe compte moins.
  Ne pas résumer le talk ici : les slides 31 à 33 portent la sortie.
---

<div class="large-space"></div>

### Le product engineer avait besoin d'un artefact pour porter l'intention.

#### On vient de passer 35 minutes à décrire cet artefact, et son outillage.

<div class="large-space"></div>

<blockquote class="charniere">
  <h4>Sans SDD, un product engineer n'est qu'un dev surchargé à qui on a retiré son PM.</h4>
  <h5>Il a besoin d'un artefact pour porter l'intention jusqu'à la machine. Cet artefact, c'est la <em>spec</em>.</h5>
</blockquote>

<div class="large-space"></div>

<div class="grid center-align">
  <div class="s2"><h6>Cadrer un problème</h6></div>
  <div class="s3"><h6>Modéliser un domaine</h6></div>
  <div class="s3"><h6>Écrire des critères vérifiables</h6></div>
  <div class="s2"><h6>Arbitrer</h6></div>
  <div class="s2"><h6>Concevoir la vérification</h6></div>
  <div class="s12"><h5><em>La vitesse de frappe compte moins.</em></h5></div>
</div>

---
layout: default
title: "Les risques, honnêtement"
speakernotes: |
  S31 - Tous les ingénieurs ne peuvent pas devenir de bons product engineers :
  le rôle exige profondeur technique, intuition produit, communication,
  conscience business et forte autogestion.
  L'erreur la plus dangereuse : déléguer l'autorité produit trop tôt, sans
  supervision suffisante ni maturité organisationnelle.
  Ajouter : atrophie des compétences chez les juniors, capacité de revue,
  responsabilité en cas d'incident.
  Rappeler l'avertissement de la product engineer française (Flowie, arrivée
  par la voie PM) : beaucoup de gens qui veulent faire du produit veulent faire
  de la stratégie et prendre des décisions, alors qu'au quotidien c'est surtout
  de la delivery. Ne devenez pas product engineer pour faire de la stratégie.
  Ce n'est pas la mort du PM : le manifeste s'adresse directement aux designers
  dans une lettre ouverte pour désamorcer l'idée d'un empiètement, et le rôle
  reste compatible avec un PM qui garde du recul sur la priorisation globale.
  Sources : CIO, The rise of the product engineer ; Le Talent Club ; Product
  Engineer Manifesto.
---

*[S31 &mdash; à produire]*

---
layout: content
title: Lundi matin
speakernotes: |
  S32 - Cinq actions par coût croissant.
  1. Écrire un AGENTS.md / constitution.md pour un repo.
  2. Utiliser le mode plan systématiquement pendant une semaine.
  3. Faire une vraie feature avec OpenSpec.
  4. Versionner les specs dans git à côté du code.
  5. Instaurer un seul gate humain - validation de spec avant implémentation -
  et mesurer le temps de revue avant/après.
  À produire : un AGENTS.md / constitution.md d'exemple, court, lisible en 20
  secondes à l'écran.
---

*[S32 &mdash; à produire]*

---
layout: default
title: La suite de la demi-journée
speakernotes: |
  S33 - Ce qui sera fait en atelier cet après-midi, et les 3 questions laissées
  ouvertes.
  Rappeler que la démo n'est pas dans ce talk : elle est dans l'atelier.
  Ressource à partager après le talk (FR) : medium.com/@damien.gouron,
  Spec-Driven Development, le guide pratique.
  Puis les coordonnées et les questions.
---

*[S33 &mdash; à produire]*
