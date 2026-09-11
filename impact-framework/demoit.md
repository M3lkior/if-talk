---
layout: cover
---

---
title: L'impact du numérique
source: https://www.arcep.fr/la-regulation/grands-dossiers-thematiques-transverses/lempreinte-environnementale-du-numerique.html
speakernotes: |
  ADEME: Agence de l'environnement et de la maîtrise de l'énergie part of Agence de la Transition Ecologique
---

<!--    <h3>D'après l'ADEME,</h3>-->

## *2,5%* des émissions de gaz à effet de serre de la France sont dues au numérique.

## *10%* de l’électricité en France est consommée par le numérique.
---
layout: default
title: Et l'IA dans tout ça ?
source: https://arxiv.org/pdf/2311.16863 (Etude  Hugging Face - Octobre 2024)
class: slide-main slide-narrow slide-prose flex flex-col items-center justify-center
---

<!--
<header>
    <nav>
        <h2 class="max center-left">Une empreinte plus large que les gaz à effet de serre</h2>
        <div class="grid center-right">
            <div class="circle transparent s6">
                <img class="responsive" src="/images/zatsit_logo.svg">
            </div>
            <div class="circle transparent s6">
                <img class="responsive" src="/images/logo.jpg">
            </div>
        </div>
    </nav>
</header>
<main class="main responsive  large-height center-align middle-align">
    <blockquote>
        <h5>
            Au-delà des gaz à effet de serre et dans un pays où la consommation énergétique est relativement
            décarbonnée,
            il est également nécessaire d’élargir la question de l’empreinte environnementale du numérique <em>à
            l’ensemble du cycle de vie des réseaux</em>,
            <em>des équipements et des terminaux</em> en adoptant une approche multicritères (terres rares, eau, énergie
            primaire…) mais également <em>leur durée de vie et les conditions de leur recyclage</em>.
        </h5>
    </blockquote>
    <br/>
</main>
<div class="center-align middle-align">
    <span>Source: https://www.arcep.fr/la-regulation/grands-dossiers-thematiques-transverses/lempreinte-environnementale-du-numerique.html</span>
</div>
-->

<blockquote class="border-main border-l-4 px-4 text-left" style="font-size: 2.0rem;">
  <h5><em>0.015</em> kWh consommé pour charger votre téléphone.</h5>
  <h5><em>0.042</em> kWh consommé par <em>ChatGPT</em> pour générer 1000 textes.</h5>
  <h5><em>0.080</em> kWh consommé par <em>votre TV</em> par heure de fonctionnement.</h5>
  <h5><em>1.0</em> kWh consommé par <em>votre frigo</em> sur une journée.</h5>
  <h5><em>2.9</em> kWh consommé par <em>Stable Diffusion</em> pour générer 1000 images.</h5>
</blockquote>
<br/>

<div class="flex items-center justify-center text-center">
  <h5>👉 https://huggingface.co/spaces/genai-impact/ecologits-calculator</h5>
</div>

---
layout: default
title: Je suis **Ludovic Dussart**
class: slide-main slide-narrow flex flex-col justify-center text-center
---

<!--
<header>
    <nav>
        <h2 class="max center-left">Vous êtes vous déjà demandé ?</h2>
        <div class="grid center-right">
            <div class="circle transparent s6">
                <img class="responsive" src="/images/zatsit_logo.svg">
            </div>
            <div class="circle transparent s6">
                <img class="responsive" src="/images/logo.jpg">
            </div>
        </div>
    </nav>
</header>
<main class="main responsive  large-height center-align middle-align">
    <blockquote>
        <h5>
            Comeibne ca représenté en banane, chocolat ou en morts ?
        </h5>
    </blockquote>
    <br/>
</main>
<div class="center-align middle-align">
    <span>Source: https://www.arcep.fr/la-regulation/grands-dossiers-thematiques-transverses/lempreinte-environnementale-du-numerique.html</span>
</div>
-->

<div class="grid grid-cols-12 items-center gap-8">
  <div class="col-span-6">
    <div>
      <img class="rounded-card" style="block-size: 15rem; " src="/images/me.jpg">
    </div>
  </div>
  <div class="col-span-6">
    <div class="text-center">
      <h4>Solutions architect <em>@zatsit</em></h4>
      <h4><em>AsyncAPI</em> maintainer</h4>
      <h4>Open Source <em>fanatic</em></h4>
    </div>
  </div>
  <div class="col-span-12"><h2 class="text-left"><em>Où me trouver ?</em></h2></div>
  <div class="contact-row"><img src="/images/twitter.jpg"/><h5>@ldussart</h5></div>
  <div class="contact-row"><img src="/images/bsky.png"/><h5>ldussart.bsky.social</h5></div>
  <div class="contact-row"><img src="/images/linkedin.png"/><h5>Ludovic Dussart</h5></div>
  <div class="col-span-12"><h2 class="text-left"><em>Où nous trouver ?</em></h2></div>
  <div class="contact-row"><img src="/images/web.png"/><h5>https://zatsit.fr</h5></div>
  <div class="contact-row"><img src="/images/bsky.png"/><h5>zatsit.bsky.social</h5></div>
  <div class="contact-row"><img src="/images/blog.png"/><h5>https://blog.zatsit.fr</h5></div>
</div>
<div class="h-8"></div>
<div>
  <h6>« Engager notre <em>expertise numérique</em> au service de <em>l'impact des entreprises</em>, en créant un écosystème <em>durable</em>, <em>partenarial</em> et <em>positif</em> »</h6>
</div>

---
layout: default-h3
title: Quels outils pour la mesure d'empreinte environnementale dans le digital ?
source: https://github.com/Green-Software-Foundation/awesome-green-software
class: slide-main slide-narrow flex items-center justify-center
speakernotes: |
  une panoplie d'outils par catégorie mais comment aggreger tout ça ?
---

<!-- beercss gave this table inline-size:100% and border-spacing:0, a bottom
     rule on every row but the last from table.border, no cell padding from
     small-space, and 1.6rem type from the talk's own xlarge-text. -->
<table class="w-full border-collapse text-[1.6rem] [&_td]:p-0 [&_th]:p-0 [&_tbody_tr:not(:last-child)_td]:border-b [&_tbody_tr:not(:last-child)_td]:border-outline">
  <tr>
    <th>Code based</th>
    <td>codecarbon.io</td>
    <td>JoularJX</td>
    <td>ecoCode</td>
    <td>etc</td>
  </tr>
  <tr>
    <th>Web</th>
    <td>ecoIndex</td>
    <td>EcoGrader</td>
    <td>Website Carbon Calculator</td>
    <td>GreenIT-Analysis</td>
    <td>etc</td>
  </tr>
  <tr>
    <th>Infra / Hardware</th>
    <td>Kepler</td>
    <td>Scaphandre</td>
    <td>PowerJoular</td>
    <td>CO2 Scope</td>
    <td>etc</td>
  </tr>
  <tr>
    <th>Cloud based</th>
    <td>Cloud Carbon Footprint</td>
    <td>Customer Carbon Footprint Tool for AWS</td>
    <td>Microsoft Emissions Impact Dashboard</td>
    <td>Carbon Footprint</td>
    <td>OVHcloud Carbon Calculator</td>
    <td>etc</td>
  </tr>
  <tr>
    <th>AI</th>
    <td>carbontracker</td>
    <td>Experiment Impact Tracker Library</td>
  </tr>
</table>

---
layout: content
title: On en construit des référentiels chez **zatsit**
class: slide-main h-stage-xlarge
---

<div class="grid h-full grid-cols-12 gap-4">
  <div class="col-span-4">
    <a class="block text-base" href="https://github.com/zatsit-oss/awesome-impact-tools"
       target="_blank"><em><strong>zatsit</strong> OSS awesome-impact-tools</em></a>
    <img class="mx-auto h-stage-large max-w-full object-contain" src="/images/zatsit-awesome.png"/>
  </div>
  <div class="col-span-8">
    <a class="block text-base" href="https://sustainability.zatsit.fr/landscape/?group=all&view-mode=grid"
       target="_blank"><em><strong>zatsit</strong> sustainability landscape</em></a>
    <web-browser src="https://sustainability-ldscp.zatsit.fr/?view-mode=grid"></web-browser>
  </div>
</div>

---
layout: content
title: Et on s'est intéressé à **Impact Framework**
class: slide-main h-stage-large
---

<div class="grid h-full grid-cols-12 gap-4">
  <div class="col-span-12 text-center">
    <h3>La mesure d'impact des logiciels doit s'appuyer sur des <em>standards</em> et des <em>outils</em>
    </h3>
  </div>
  <div class="col-span-12 w-full text-center">
    <img class="mx-auto max-h-full object-contain" src="/images/gsf.avif"/>
  </div>
</div>

---
layout: content
title: Motivations
class: slide-main h-stage-large text-center
speakernotes: |
  La mesure de l'impact des logiciels sur des paramètres tels que le carbone, l'eau et l'énergie est
  complexe et nuancée.
  <br/><br/>
  Les applications modernes sont composées de nombreux éléments logiciels plus petits (composants)
  fonctionnant sur différents environnements, par exemple, cloud privé, cloud public, bare-metal, virtualisé,
  conteneurisé, mobile, ordinateurs portables, ordinateurs de bureau, embarqué et IoT.
  <br/><br/>
  Transformer des observations en impacts: Les mesures facilement observables telles que l'utilisation de l'unité
  centrale, les pages vues et les installations sont converties en impacts environnementaux tels que les émissions de
  carbone, l'utilisation de l'eau, la consommation d'énergie et la qualité de l'air.<br/><br/>
  Explorer les scénarios de simulation: Découvrez l'impact environnemental des changements apportés à votre logiciel
  en temps réel en examinant comment la performance environnementale de votre logiciel évolue si votre application est
  déplacée vers le nuage ou si son temps d'exécution est modifié.<br/><br/>
  Stockez:
  Enregistrez vos observations, les plugins choisis, les configurations et les impacts environnementaux calculés dans
  un fichier manifeste. Ouvrez la porte à d'autres personnes pour qu'elles comprennent, vérifient et remettent en
  question l'ensemble du processus.
  Permettez à d'autres personnes de réexécuter votre fichier manifeste, de valider vos résultats ou de remettre en
  question vos hypothèses. Ils peuvent modifier les configurations, sélectionner différents plugins et effectuer
  eux-mêmes l'analyse.
---

<div class="h-12"></div>
<h3>Impact Framework (IF) vise à faciliter le <em>calcul</em> et le <em>partage</em> des impacts
environnementaux des logiciels.</h3>

<div class="h-12"></div>
<div class="h-12"></div>
<h3 class="text-left">Les promesses du framework ?</h3>
<div class="h-12"></div>
<div class="h-12"></div>
<div class="absolute left-1/2 -translate-x-1/2">
  <h4 class="text-left"><em>Transformez</em> vos observations en unités d'impacts.</h4>
  <h4 class="text-left"><em>Explorez</em> des scénarios de simulation.</h4>
  <h4 class="text-left"><em>Stockez</em> vos manifestes et <em>démocratisez-</em>les.</h4>
</div>

---
layout: split
title: IF - Concepts majeurs
height: xlarge
speakernotes: |
  - D'abord présenter if-run
  - Ensuite montrer le manifest
---

::browser{src=https://if.greensoftware.foundation/major-concepts/if}

---
layout: split
title: IF - Fonctionnement des pipelines
height: xlarge
---

```mermaid
                block-beta
                    columns 1
                        plugin1["Mock d'observations"]
                        space
                        block:out
                            cpu_out[/"cpu/utilization"/]
                            memory_out[/"memory/utilization"/]
                            memory_capacity_out[/"memory/capacity"/]
                        end
                        space
                        plugin2["Consommation d'énergie par la mémoire"]
                        space
                        memory_energy[/"memory/energy"/]
                        space
                        plugin3["Et ainsi de suite ..."]

                    plugin1 -- "output" --> cpu_out
                    plugin1 -- "output" --> memory_out
                    memory_capacity_out -- "input" --> plugin2
                    memory_out -- "input" --> plugin2
                    plugin2 -- "output" --> memory_energy
                    memory_energy -- "input" --> plugin3
                    classDef Input fill:#0f15fd,stroke:white,color:white,border:white;
                    %%classDef Output fill:#f1be51,stroke:white;
                    classDef Plugin stroke:#0f15fd,color:#0f15fd,stroke-width:2px,stroke-dasharray: 5 5

                    class cpu_out,memory_out,memory_capacity_out,memory_energy Input
                    class plugin1,plugin2,plugin3 Plugin
                    style out fill:#fff,stroke:white;
```

::code{folder=sources files=pipelines.yml lines=11-20}

---
layout: split
title: IF - Hello World !
height: xlarge
speakernotes: |
  - montrer if-run --help
  - energy : kwH
---

::term{path=sandbox}

::browser{src=https://if.greensoftware.foundation/users/quick-start}

---
layout: content
title: IF - Observer.
speakernotes: |
  - montrer if-run --help
  - energy : kwH
  - carbon : eqCo2
---

:::grid{class="grid grid-cols-12 grid-rows-[auto_1fr] gap-4 h-stage-xlarge"}
:::col{class="col-span-12"}
#### Pour pouvoir commencer à *mesurer*, il faut d'abord *observer* et *récolter* de la data.
:::

:::col{class="col-span-4 h-stage-xlarge"}
<div>

<br/>
<p class="text-[1.6rem]">Il existe quelques plugins permettant de nourrir les pipelines en <em>inputs
:</em>
</p>

<!-- list-disc and pl-8 are asked for explicitly: beercss restored native list
     markers with `:not(nav) > :is(ul,ol) { all: revert }`, and Tailwind's
     preflight does not. -->
<ul class="list-disc pl-8 text-[1.6rem]">
  <li>Mock Observations</li>
  <li>AWS Importer</li>
  <li>Azure Importer</li>
  <li>Datadog Importer</li>
  <li>cloud-storage-metadata</li>
  <li>prometheus-importer</li>
  <li>et pleins d'autres à suivres ...</li>
</ul>

</div>
:::

:::col{class="col-span-8 h-stage-xlarge"}
::vscode{path=sources}

<!--        <web-browser src="https://explorer.if.greensoftware.foundation/"></web-browser>-->
:::
:::

---
layout: content
title: IF - Observer
speakernotes: |
  - montrer group
  - puis cloud-metadata
---

:::split{cols=4,8 height=xlarge}
::term{path=sources}

::vscode{path=sources}
:::

---
layout: content
title: IF - C'est quoi la consommation energétique de zatsit.fr ?
class: slide-main h-stage-xlarge
---

<!--
<header>
    <nav>
        <h5 class="max center-left">IF - Calculer.</h5>
        <div class="grid center-right">
            <div class="circle transparent s6">
                <img class="responsive" src="/images/zatsit_logo.svg">
            </div>
            <div class="circle transparent s6">
                <img class="responsive" src="/images/logo.jpg">
            </div>
        </div>
    </nav>
</header>
<main class="responsive max xlarge-height">
    <text>
        <h5>Cela consiste à <em>assembler</em> les plugins pour obtenir des <em>outputs</em> de type
            <em>mesures</em>.</h5>
        <br/>
    </text>
    <split-view class="large-height">
        <img src="/images/if-plugins.png" style="width: 100%"/>
        &lt;!&ndash;            <img src="/images/if-unofficial.png"/>&ndash;&gt;
        <web-browser src="https://explorer.if.greensoftware.foundation/"></web-browser>
    </split-view>
    <speaker-notes>
        - montrer if-run &#45;&#45;help
        - energy : kwH
        - carbon : eqCo2
    </speaker-notes>
</main>
-->

:::split{cols=4,8 height=xlarge}
::term{path=sources}

::vscode{path=sources}
:::

---
layout: content
title: IF - Visuellement, c'est mieux.
class: slide-main h-stage-xlarge
---

:::split{cols=4,8 height=xlarge}
::term{path=sources}

::vscode{path=sources}
:::

---
layout: content
title: IF - Visuellement, c'est mieux.
class: slide-main h-stage-xlarge
---

:::split{cols=4,8 height=xlarge}
::term{path=sources}

::browser{src=http://localhost:3000}
:::

---
layout: content
title: IF - C'est quoi la suite ?
class: slide-main flex flex-col items-center justify-center gap-12 text-center
---

<h3>Imaginez la configuration des pipelines Impact Framework (IF) en WYSIWYG ?</h3>

<div class="text-left">
  <h4><em>Essayez</em> le framework.</h4>
  <h4><em>Contribuer</em> et faire évoluer l'existant.</h4>
  <h4><em>Proposez</em> de nouveaux plugins ?</h4>
  <h4><em>Constituez</em> vos propres manifests.</h4>
</div>

<img src="/images/hacktoberfest.png" style="width: 20rem"/>

---
layout: default
title: Merci de votre attention. Des questions ?
class: slide-main slide-narrow flex flex-col justify-center gap-8 text-center
---

<div class="w-full">
  <img src="/images/qa.png" style="block-size: 22rem;"/>
</div>

<div class="grid grid-cols-12 items-center gap-8">
  <div class="col-span-12"><h2 class="text-left"><em>Où me trouver ?</em></h2></div>
  <div class="contact-row"><img src="/images/twitter.jpg"/><h5>@ldussart</h5></div>
  <div class="contact-row"><img src="/images/bsky.png"/><h5>ldussart.bsky.social</h5></div>
  <div class="contact-row"><img src="/images/linkedin.png"/><h5>Ludovic Dussart</h5></div>
  <div class="col-span-12"><h2 class="text-left"><em>Où nous trouver ?</em></h2></div>
  <div class="contact-row"><img src="/images/web.png"/><h5>https://zatsit.fr</h5></div>
  <div class="contact-row"><img src="/images/bsky.png"/><h5>zatsit.bsky.social</h5></div>
  <div class="contact-row"><img src="/images/blog.png"/><h5>https://blog.zatsit.fr</h5></div>
</div>
