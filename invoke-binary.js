/**
 * Javascript shim inspired by from Blend:
 * https://full-stack.blend.com/how-we-write-github-actions-in-go/invoke-binary.js
 */

const childProcess = require('child_process')
const os = require('os')
const process = require('process')

const VERSION = '"__VERSION_SHA__"'

function chooseBinary() {
    const platform = os.platform()
    const arch = os.arch()

    const platformMap = {
        linux: 'linux',
        windows: 'windows',
    };
    
    const archMap = {
        x64: 'amd64',
        arm64: 'arm64',
    };
    
    if (platformMap[platform] && archMap[arch]) {
        return `action-${platformMap[platform]}-${archMap[arch]}-${VERSION}`;
    }

    console.error(`unsupported platform (${platform}) and architecture (${arch})`)
    process.exit(1)
}

function main() {
    const binary = chooseBinary()
    const mainScript = `${__dirname}/${binary}`
    const spawnSyncReturns = childProcess.spawnSync(mainScript, { stdio: 'inherit' })
    const status = spawnSyncReturns.status
    if (typeof status === 'number') {
        process.exit(status)
    }
    process.exit(1)
}

if (require.main === module) {
    main()
}