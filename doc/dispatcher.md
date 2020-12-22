```dot
digraph {
    rankdir=LR

    i [label=logger]
    d [label=dispatcher shape=box]
    i -> d [label=item]

    subgraph cluster1 {
        label = "level writers"
        color = transparent

        w1 [shape=box label="writer\n(default)"]
        o1 [label=output]
        w1 -> o1

        w2 [shape=box style=dashed label="writer\n(additional)"]
        o2 [label=output style=dashed]
        w2 -> o2 [style=dashed]
    }    

    d -> {w1 w2} [label=item]

    d -> w1 [color=gray style=dashed label="feedback item"]
    w2 -> d [color=gray style=dashed label="feedback item"]
}
