import * as vscode from 'vscode';
import * as Parser from 'tree-sitter';
import * as FunLanguage from 'tree-sitter-fun';

export function activate(context: vscode.ExtensionContext) {
    console.log('Fun Language extension is now active!');

    // Register the tree-sitter grammar
    const funLanguage = FunLanguage;
    const parser = new Parser();
    parser.setLanguage(funLanguage);

    // Register document symbol provider for outline
    const documentSymbolProvider = vscode.languages.registerDocumentSymbolProvider(
        { language: 'fun' },
        new FunDocumentSymbolProvider(parser)
    );

    context.subscriptions.push(documentSymbolProvider);
}

export function deactivate() {}

class FunDocumentSymbolProvider implements vscode.DocumentSymbolProvider {
    constructor(private parser: Parser) {}

    provideDocumentSymbols(
        document: vscode.TextDocument,
        token: vscode.CancellationToken
    ): vscode.ProviderResult<vscode.SymbolInformation[] | vscode.DocumentSymbol[]> {
        const text = document.getText();
        const tree = this.parser.parse(text);
        const symbols: vscode.DocumentSymbol[] = [];

        this.extractSymbols(tree.rootNode, document, symbols);
        return symbols;
    }

    private extractSymbols(
        node: Parser.SyntaxNode,
        document: vscode.TextDocument,
        symbols: vscode.DocumentSymbol[]
    ) {
        // Extract variable assignments
        if (node.type === 'assign') {
            const children = node.children;
            if (children.length >= 3 && children[0].type === 'var') {
                const nameNode = children[0];
                const range = new vscode.Range(
                    document.positionAt(nameNode.startIndex),
                    document.positionAt(nameNode.endIndex)
                );
                const symbol = new vscode.DocumentSymbol(
                    nameNode.text,
                    'Variable',
                    vscode.SymbolKind.Variable,
                    range,
                    range
                );
                symbols.push(symbol);
            }
        }

        // Extract function definitions (lambda expressions)
        if (node.type === 'lam') {
            const children = node.children;
            if (children.length >= 3 && children[0].type === '\\') {
                // Find parameter names
                const params: string[] = [];
                for (let i = 1; i < children.length - 2; i++) {
                    if (children[i].type === 'var') {
                        params.push(children[i].text);
                    }
                }
                const funcName = params.length > 0 ? params.join(', ') : 'λ';
                const range = new vscode.Range(
                    document.positionAt(node.startIndex),
                    document.positionAt(node.endIndex)
                );
                const symbol = new vscode.DocumentSymbol(
                    funcName,
                    'Function',
                    vscode.SymbolKind.Function,
                    range,
                    range
                );
                symbols.push(symbol);
            }
        }

        // Extract constructors
        if (node.type === 'cons') {
            const consNameNode = node.children.find(child => child.type === 'cons_name');
            if (consNameNode) {
                const range = new vscode.Range(
                    document.positionAt(node.startIndex),
                    document.positionAt(node.endIndex)
                );
                const symbol = new vscode.DocumentSymbol(
                    consNameNode.text,
                    'Constructor',
                    vscode.SymbolKind.Class,
                    range,
                    range
                );
                symbols.push(symbol);
            }
        }

        // Extract record definitions
        if (node.type === 'rec') {
            const range = new vscode.Range(
                document.positionAt(node.startIndex),
                document.positionAt(node.endIndex)
            );
            const symbol = new vscode.DocumentSymbol(
                'Record',
                'Record',
                vscode.SymbolKind.Struct,
                range,
                range
            );
            symbols.push(symbol);
        }

        // Extract record properties
        if (node.type === 'prop') {
            const children = node.children;
            if (children.length >= 3 && children[2].type === 'var') {
                const propNameNode = children[2];
                const range = new vscode.Range(
                    document.positionAt(node.startIndex),
                    document.positionAt(node.endIndex)
                );
                const symbol = new vscode.DocumentSymbol(
                    propNameNode.text,
                    'Property',
                    vscode.SymbolKind.Property,
                    range,
                    range
                );
                symbols.push(symbol);
            }
        }

        // Extract list literals
        if (node.type === 'list') {
            const range = new vscode.Range(
                document.positionAt(node.startIndex),
                document.positionAt(node.endIndex)
            );
            const symbol = new vscode.DocumentSymbol(
                'List',
                'List',
                vscode.SymbolKind.Array,
                range,
                range
            );
            symbols.push(symbol);
        }

        // Extract string literals
        if (node.type === 'str') {
            const range = new vscode.Range(
                document.positionAt(node.startIndex),
                document.positionAt(node.endIndex)
            );
            const symbol = new vscode.DocumentSymbol(
                'String',
                'String',
                vscode.SymbolKind.String,
                range,
                range
            );
            symbols.push(symbol);
        }

        // Extract integer literals
        if (node.type === 'int') {
            const range = new vscode.Range(
                document.positionAt(node.startIndex),
                document.positionAt(node.endIndex)
            );
            const symbol = new vscode.DocumentSymbol(
                node.text,
                'Integer',
                vscode.SymbolKind.Number,
                range,
                range
            );
            symbols.push(symbol);
        }

        // Extract pattern matching (when expressions)
        if (node.type === 'when') {
            const range = new vscode.Range(
                document.positionAt(node.startIndex),
                document.positionAt(node.endIndex)
            );
            const symbol = new vscode.DocumentSymbol(
                'Pattern Match',
                'Pattern Match',
                vscode.SymbolKind.Method,
                range,
                range
            );
            symbols.push(symbol);
        }

        // Extract function applications
        if (node.type === 'app') {
            const range = new vscode.Range(
                document.positionAt(node.startIndex),
                document.positionAt(node.endIndex)
            );
            const symbol = new vscode.DocumentSymbol(
                'Function Call',
                'Function Call',
                vscode.SymbolKind.Method,
                range,
                range
            );
            symbols.push(symbol);
        }

        // Extract infix applications
        if (node.type === 'iapp') {
            const children = node.children;
            if (children.length >= 3 && children[1].type === 'sym') {
                const operatorNode = children[1];
                const range = new vscode.Range(
                    document.positionAt(node.startIndex),
                    document.positionAt(node.endIndex)
                );
                const symbol = new vscode.DocumentSymbol(
                    operatorNode.text,
                    'Infix Operation',
                    vscode.SymbolKind.Operator,
                    range,
                    range
                );
                symbols.push(symbol);
            }
        }

        // Extract blocks
        if (node.type === 'block') {
            const range = new vscode.Range(
                document.positionAt(node.startIndex),
                document.positionAt(node.endIndex)
            );
            const symbol = new vscode.DocumentSymbol(
                'Block',
                'Block',
                vscode.SymbolKind.Namespace,
                range,
                range
            );
            symbols.push(symbol);
        }

        // Extract type annotations
        if (node.type === 'annot') {
            const children = node.children;
            if (children.length >= 3 && children[0].type === 'var') {
                const nameNode = children[0];
                const range = new vscode.Range(
                    document.positionAt(node.startIndex),
                    document.positionAt(node.endIndex)
                );
                const symbol = new vscode.DocumentSymbol(
                    nameNode.text,
                    'Type Annotation',
                    vscode.SymbolKind.TypeParameter,
                    range,
                    range
                );
                symbols.push(symbol);
            }
        }

        // Extract imports
        if (node.type === 'import') {
            const children = node.children;
            if (children.length >= 5 && children[1].type === 'var') {
                const nameNode = children[1];
                const range = new vscode.Range(
                    document.positionAt(node.startIndex),
                    document.positionAt(node.endIndex)
                );
                const symbol = new vscode.DocumentSymbol(
                    `import ${nameNode.text}`,
                    'Import',
                    vscode.SymbolKind.Module,
                    range,
                    range
                );
                symbols.push(symbol);
            }
        }

        // Recursively process child nodes
        for (const child of node.children) {
            this.extractSymbols(child, document, symbols);
        }
    }
} 