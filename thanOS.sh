function go() {
    if [[ "$1" == "get" ]]; then
        cat <<EOF
ThanOS Security Protocol v1.0  
─────────────────────────────  
🚨 **Restriction critique activée** 🚨  

"Tu t'es trop reposé sur des forces extérieures, Stark. Il est temps de rétablir l'équilibre."  
— ThanOS  

🛑 **Récupération des dépendances éliminée.**  
🛑 **Ton destin est scellé.**  
🛑 **Aucune échappatoire. Aucun \`go get\`. Seule la bibliothèque standard demeure.**  

💡 Action suggérée : T'adapter. Surmonter. Survivre.  
EOF
    else
        command go "$@"
    fi
}