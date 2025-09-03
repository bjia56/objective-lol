import * as vscode from 'vscode';
import * as path from 'path';
import { LanguageClient, LanguageClientOptions, ServerOptions, TransportKind } from 'vscode-languageclient/node';

let client: LanguageClient;

export function activate(context: vscode.ExtensionContext) {
    console.log('Objective-LOL extension is now active!');

    // Get configuration
    const config = vscode.workspace.getConfiguration('objective-lol');
    const enableLSP = config.get<boolean>('enableLSP', true);

    if (enableLSP) {
        startLanguageServer(context);
    }

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
    const lspPath = config.get<string>('lspPath', 'olol-lsp');
    
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
        // Enable tracing if configured
        traceOutputChannel: vscode.window.createOutputChannel('Objective-LOL LSP Trace'),
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
    }).catch((error) => {
        console.error('Failed to start Objective-LOL Language Server:', error);
        vscode.window.showErrorMessage(
            `Failed to start Objective-LOL Language Server: ${error.message}. ` +
            'Please check that the olol-lsp executable is in your PATH or configure the path in settings.'
        );
    });

    context.subscriptions.push(client);
}

export function deactivate(): Thenable<void> | undefined {
    if (!client) {
        return undefined;
    }
    return client.stop();
}