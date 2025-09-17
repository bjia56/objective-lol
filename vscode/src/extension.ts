// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as vscode from 'vscode';
import * as path from 'path';
import * as fs from 'fs';
import { LanguageClient, LanguageClientOptions, ServerOptions } from 'vscode-languageclient/node';

let client: LanguageClient;

// This method is called when your extension is activated
// Your extension is activated the very first time the command is executed
export function activate(context: vscode.ExtensionContext) {
    console.log('Objective-LOL extension is now active!');

    ensureLanguageServerInstalled(context).then(() => {
        startLanguageServer(context);
    }).catch((error) => {
        console.error('Error ensuring language server is installed:', error);
    }).then(() => {
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
    });
}

const lspVersion = 'v0.0.2';
const lspPlatform = (() => {
    switch (process.platform) {
        case 'win32':
            return 'windows';
        default:
            return process.platform;
    }
})();
const lspArch = (() => {
    switch (process.arch) {
        case 'x64':
            return 'amd64';
        default:
            return process.arch;
    }
})();

const suffix = lspPlatform === 'windows' ? '.exe' : '';
const binaryName = `olol-lsp-${lspVersion}-${lspPlatform}-${lspArch}${suffix}`;

function ensureLanguageServerInstalled(context: vscode.ExtensionContext) {
    const url = `https://github.com/bjia56/objective-lol/releases/download/${lspVersion}/olol-lsp-${lspPlatform}-${lspArch}${suffix}`;

    const binaryPath = path.join(context.globalStorageUri.fsPath, binaryName);
    console.log(`Checking for Objective-LOL Language Server at ${binaryPath}`);

    return new Promise<void>((resolve, reject) => {
        if (!fs.existsSync(binaryPath)) {
            vscode.window.showInformationMessage('Downloading Objective-LOL Language Server...');
            fs.mkdirSync(context.globalStorageUri.fsPath, { recursive: true });

            fetch(url).then(res => {
                if (res.status !== 200) {
                    reject(new Error(`Failed to download language server: ${res.statusText}`));
                    return;
                }

                res.arrayBuffer().then(buffer => {
                    fs.writeFileSync(binaryPath, Buffer.from(buffer), { mode: 0o755 });
                    vscode.window.showInformationMessage('Objective-LOL Language Server downloaded successfully.');
                    console.log('Objective-LOL Language Server downloaded and installed.');
                    resolve();
                }).catch(err => {
                    reject(new Error(`Error reading response buffer: ${err.message}`));
                });
            }).catch(err => {
                reject(new Error(`Error fetching language server: ${err.message}`));
            });
        } else {
            console.log('Objective-LOL Language Server already installed.');
            resolve();
        }
    });
}

function startLanguageServer(context: vscode.ExtensionContext) {
    const lspPath = path.join(context.globalStorageUri.fsPath, binaryName);

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
            `Failed to start Objective-LOL Language Server: ${error.message}. `
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
