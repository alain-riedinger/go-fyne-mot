# go-fyne-mot

Le jeu du **Mot le Plus Long** est un jeu de lettres issu du jeu télévisé "Des Chiffres et des Lettres": à partir d'un tirage des 10 lettres, voyelles et consonnes, il faut trouver le plus long mot possible. Les variations suivantes autour d'un mot de base sont acceptées:  
- pour les verbes: l'infinitif, le participe passé, le participe présent
- pour les noms communs: le singulier et le pluriel
- pour les adjectifs: le masculin, le féminin et le pluriel, ainsi que toute les combinaisons possibles

Les noms propres ne sont pas acceptés.

## Le dictionnaire

Le programme s'appuye sur un dictionnaire, basé sur **hunspell**, c'est ce dictionnaire qui est utilisé dans les principaux programmes *Open Source*, comme "LibreOffice" ou "Firefox".  

Pour des raisons de praticité, les manipulations ont été effectuées sur un **Raspberry Pi 4**, sous **Raspbian**: tous les outils nécessaires sont disponibles et les manipulations bien décrites sur le Net.
- installation de **hunspell**
```bash
sudo apt install hunspell
sudo apt install hunspell-fr-comprehensive
```
- installation des outils annexes, dont **unmunch** qui est nécessaire pour *aplanir* le dictionnaire
```bash
sudo apt install hunspell-tools
```

Les dictionnaires hunspell pour une langue donnée sont composés de 2 fichiers:  
- le dictionnaire, qui contient tous les mots *racine*:  
"/usr/share/hunspell/fr_FR.dic"
- les affixes, qui contient toutes les variations des mots *racine* (féminin, pluriel, conjugaisons):  
"/usr/share/hunspell/fr_FR.aff"

Pour obtenir un dictionnaire *aplani* compatible avec les règles du **Mot le Plus Long**, il faut créer un fichier avec des affixes épuré des combinaisons interdites (Cf les règles plus haut): "fr-mlpl.aff". Ce fichier est édité à la main et modifié pour ne garder que les conjugaisons autorisées. C'est assez fastidieux, mais j'ai opté pour une vérification pour chaque bloc, parce que **unmunch** est très peu verbeux en cas d'erreur.
```bash
unmunch /usr/share/hunspell/fr_FR.dic /home/pi/Documents/fr-mlpl.aff > fr-mlpl-flat.txt
```

En cas d'erreur, le bloc modifié est restauré, puis corrigé en faisant bien attention à chaque modification.  
Les conjugaisons commencent à la ligne **702** et pour chaque bloc, regroupant les verbes de même racine (comme *prendre* et *surprendre*) sont regroupés dans un même bloc.  
Il faut effectuer les opérations suivantes sur chaque bloc:
- supprimer les conjugaisons autres que: infinitif, participe présent, participe passé
- compter le nombre de lignes restantes
- mettre à jour l'entête du bloc, avec le nombre exact de lignes dans le bloc, par exemple, s'il reste **9** lignes dans le bloc:
```
SFX zA Y 9
```
Le fichier d'affixes est également nettoyé des unité abrégées, comme "k" au lieu de "kilo".  
ligne 652
```
PFX U. Y 20
PFX U. 0 Y .
PFX U. 0 Z .
...
PFX U. 0 y .
```

Ce dictionnaire *aplani* nommé ici "fr-mlpl-flat.txt" va servir de base à la réalisation du dictionnaire dédié au **Mot le Plus Long**.  

Le programme dispose de 2 modes:
- le mode *création du dictionnaire*, pour créer le dictionnaire dédié au jeu
- le mode de jeu lui-même

Pour trouver les numéros de ligne définissant l'intervalle, il faut chercher ces 2 lignes:
> a  
> ...  
> zzzz  

Il faut exécuter **une seule fois** la création du dictionnaire avec la commande:
```
go-MotLePlusLong.exe dico fr-mlpl-flat.txt --start 26343 --end 806546
```

Un fichier "fr-mlpl-flat-strict.txt" est généré, il contient le dictionnaire *aplani* avec les modifications suivantes, pour ne garder que les mots autorisés par les règles du jeu:
- suppression de toutes les suggestions propres au format de hunspell: "s'...", "blabla/..."
- suppression de tous *les noms propres*
- suppression de tous les mots *de plus de 10 lettres*
- suppression de toutes les expression (avec des quotes "'")
- suppression des mots avec des majuscules
- retrait du caractère "-" des mots composés
- remplacement de tous les caractères étendus par leur caractère de base de l'alphabet
par exemple: **éèëê** deviennent un **e**, donc "mètre" devient "metre"

Une passe manuelle est nécessaire pour supprimer tous les mots avec des unités de mesures avec les préfixes en abrégé, mais en gardant les préfixes en toutes lettres:
- "kilometre" est accepté
- "kmetre" ou "km" sont supprimés

Le fichier "fr-mlpl-flat-strict.txt" finalement obtenu est la référence même pour le jeu.  
Il contient un total de **111496**.  
A noter: je ne garantis pas que ce dictionnaire *fait maison* soit conforme avec les 2 dictionnaires de référence utilisés dans le jeu officiel (le *Larousse* et le *Petit Robert*).

## Le programme

Le programme en lui-même est fait en **Go**, c'est une occasion d'apprendre ce langage, avec des challenges:
- gérer la récursivité
- calculer des hashs
- manipuler des dictionnaires et des tableaux

Lors de la conversion des lettres accentuées en lettres simples, il a fallu composer avec les particularités de Go: il faut basculer les `string` en `rune` pour que les accents soient gérés correctement:
```
	// Go UTF8 and Unicode interaction needs a string and rune equivalent:
	// à   â   ä   é   è   ê   ë   î   ï   ô   ö   ù   û   ü   ç
	// 224 226 228 233 232 234 235 238 239 244 246 249 251 252 231
	rline := []rune(line)
```

Grace au dictionnaire *aplani*, l'algorithme de solution du jeu est devenu simple:
- le programme charge au lancement le dictionnaire
- chaque mot est indexé en fonction du nombre total de lettres contenu, puis des nombres de chaque lettre contenue, pour les 26 lettres de a à z. Le tout est encodé sur un `[14]byte` de **14 octets**. Le premier octet contient la longueur du mot, les 13 suivants le nombre d'occurence de chaque lettre, encodé par deux sur un même octet.  

Pour chaque index, il y a une liste de mots *similaires* (qui ont le même index, en bref des anagrammes), le tout est stocké dans un `dico map[[14]byte][]string`.  
L'algorithme consiste donc à chercher les mots les plus longs, en commençant par les 10 lettres (1 seule possibilité), puis s'il n'y a pas de correspondance, tentative avec les mots de 9 lettres, en retirant à tour de rôle une lettre (10 possibilités), puis même chose avec les 8 lettres, et ainsi de suite...  
Toutes les solutions trouvées pour une longueur donnée sont conservées et affichées dans la solution à la fin de la manche.  

Pour l'IHM, le framework utilisé est [fyne](https://github.com/fyne-io/fyne), dans sa version **v2** (no compatible avec les versions v1.x).  
Ce framework est suffisant pour faire une IHM minimaliste qui permet au jeu de fonctionner correctement.  
