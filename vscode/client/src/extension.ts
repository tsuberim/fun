/* --------------------------------------------------------------------------------------------
 * Copyright (c) Microsoft Corporation. All rights reserved.
 * Licensed under the MIT License. See License.txt in the project root for license information.
 * ------------------------------------------------------------------------------------------ */

// import * as path from 'path';
import { ExtensionContext } from 'vscode';

import {
	Executable,
	LanguageClient,
	LanguageClientOptions,
	TransportKind
} from 'vscode-languageclient/node';

let client: LanguageClient;

export function activate(_: ExtensionContext) {
	const serverOptions: Executable = {
		transport: TransportKind.stdio, 
		command: '/Users/mtsuberi/Projects/fun/fun',
		args: ['lsp'] 
	};

	const clientOptions: LanguageClientOptions = {
		documentSelector: [{ scheme: 'file', language: 'fun' }]
	};

	// Create the language client and start the client.
	client = new LanguageClient(
		'funLanguageServer',
		'Fun Language Server',
		serverOptions,
		clientOptions
	);

	// Start the client. This will also launch the server
	client.start();
}

export function deactivate(): Thenable<void> | undefined {
	if (!client) {
		return undefined;
	}
	return client.stop();
}
