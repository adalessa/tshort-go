# Tshort

Applicacion que permite un rapido switch entre projectos
sea creando una nueva session o moviendose

# instalacion

```
go install github.com/adalessa/tshort@latest
```

# Uso
En tmux agregar
```
bind-key S run-shell 'tmux popup -E tshort'

bind-key u run-shell  'tmux popup -E tshort switch 1'
bind-key i run-shell  'tmux popup -E tshort switch 2'
bind-key o run-shell  'tmux popup -E tshort switch 3'
bind-key p run-shell  'tmux popup -E tshort switch 4'
```

Para agregar a la barra de status utilizar el comando `tshort list`
