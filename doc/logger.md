```dot
digraph {
    rankdir=LR

    i [label=app]
    
    l [label=logger shape=box]
    d [label=dispatcher shape=box]
    
    i -> l [label=methods]
    l -> d [label=items]

    subgraph cluster1 {
        label = "logging levels"
        color = transparent

        w1 [shape=box label="[]writer"]
        o1 [label="[]output"]
        w1 -> o1

        w2 [shape=box style=dashed label="[]writer"]
        o2 [label="[]output" style=dashed]
        w2 -> o2 [style=dashed]

        w3 [shape=box label="[]writer"]
        o3 [label="[]output"]
        w3 -> o3
    }    

    d -> {w1, w3} [label=item]
    d -> {w2} [label=item style=dashed]
}
```
