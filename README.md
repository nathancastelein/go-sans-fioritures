# Go sans fioritures

Ce repository a pour but d'illustrer une conférence nommée "Go sans fioritures: quand le standard suffit".

La conférence a pour but de montrer les possibilités offertes par la librairie standard pour construire une API, en utilisant les packages `net/http`, `log/slog` et `testing`.

## Résumé de la conférence

Dès sa première version, Go a été conçu avec une philosophie minimaliste. "Faire beaucoup avec peu" a toujours été une motivation du langage et de ses librairies standard. 

Face à certains manques dans ces dernières, la communauté a fourni de nombreux outils: des framework web (gin, echo, fiber, …), des librairies de logs (logrus, zerolog) ou même divers outils de testing comme testify.

Saviez-vous qu'une large partie des fonctionnalités de ces librairies tierces avait été intégrée dans les packages standards au fil du temps ? Il est désormais possible d'écrire des applications modernes et performantes tout en réduisant la dépendance aux bibliothèques tierces.

Dans cette session, nous explorerons comment tirer le meilleur parti des packages natifs de Go, tels que net/http, log/slog, testing, embed, etc.

Nous comparerons ces fonctionnalités aux bibliothèques tierces bien connues, et montrerons pourquoi et comment revenir aux fondamentaux peut simplifier vos projets, réduire la maintenance et améliorer la compatibilité à long terme tout en embrassant la philosophie du langage Go.