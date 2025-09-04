// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as vscode from 'vscode';
import { LanguageClient, LanguageClientOptions, ServerOptions } from 'vscode-languageclient/node';

let client: LanguageClient;

// This method is called when your extension is activated
// Your extension is activated the very first time the command is executed
export function activate(context: vscode.ExtensionContext) {
    console.log('Objective-LOL extension is now active!');

    startLanguageServer(context);

    // Register commands
    context.subscriptions.push(
        vscode.commands.registerCommand('objective-lol.restartLanguageServer', () => {
            if (client) {
                client.stop().then(() => {
                    startLanguageServer(context);
                });
            } else {
                startLanguageServer(context);
            }
        })
    );
}

function startLanguageServer(context: vscode.ExtensionContext) {
    const config = vscode.workspace.getConfiguration('objective-lol');
    const lspPath = config.get<string>('lsp-path', 'olol-lsp');

    // Server options - runs the LSP server
    const serverOptions: ServerOptions = {
        command: lspPath,
        args: [],
        options: {
            shell: true
        }
    };

    // Client options - defines which files to monitor and language features to enable
    const clientOptions: LanguageClientOptions = {
        documentSelector: [
            {
                scheme: 'file',
                language: 'objective-lol'
            }
        ],
        synchronize: {
            // Synchronize the setting section 'objective-lol' to the server
            configurationSection: 'objective-lol',
            // Notify the server about file changes to '.olol' files contained in the workspace
            fileEvents: vscode.workspace.createFileSystemWatcher('**/*.olol')
        },
        outputChannelName: 'Objective-LOL Language Server'
    };

    // Create the language client
    client = new LanguageClient(
        'objective-lol-lsp',
        'Objective-LOL Language Server',
        serverOptions,
        clientOptions
    );

    // Start the client and server
    client.start().then(() => {
        console.log('Objective-LOL Language Server started successfully');
    }).catch((error: { message: any; }) => {
        console.error('Failed to start Objective-LOL Language Server:', error);
        vscode.window.showErrorMessage(
            `Failed to start Objective-LOL Language Server: ${error.message}. ` +
            'Please check that the olol-lsp executable is in your PATH or configure the path in settings.'
        );
    });

    context.subscriptions.push(client);
}

// This method is called when your extension is deactivated
export function deactivate() {
	if (!client) {
		return;
	}
	return client.stop();
}
