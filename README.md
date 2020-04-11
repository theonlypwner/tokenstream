# tokenstream
This project generates a string of non-terminal symbols (tokens) that parses to a given non-terminal in a context-free grammar.

* Parser: (tokens, grammar) -> tree
* Unparser: (tree, grammar) -> tokens
* this project: (non-terminal, grammar) -> tokens

Parsers take a string of tokens (and a grammar) and output a syntax or parse tree. Unparsers do the reverse and take a tree (and a grammar) and output a string of tokens through a simple walk through the tree.

This project takes a non-terminal (and a grammar) and directly outputs a string of terminals, saving memory by not building the entire tree.
