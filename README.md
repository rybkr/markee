# Markee
## Markee
### Markee
#### Markee
##### Markee
###### Markee

~~~py
def
~~~

> Quote
> > Quote
> > > Quote
> > Quote
> Quote

> Quote

> # Why

1. # H
2. Two
3. Three

* One
2. Two
3. Three

***

```
Each node type will have a continue() method that checks whether a line is a valid continuation of sed node type
Custom nodes will have continue as a member
Different nodes have different priorities:

this is something
this is other thing.
1. this is a list

makes a paragraph not a list.

consume continuation markers: '>' gets eaten for a blockquote.

Remaining text is chucked as the member of the last open block.

Each node type also needs a finalixe() method thsat can be a member for custom nodes.
finalize takes care of tidying up the AST before inlin eparsing


```
