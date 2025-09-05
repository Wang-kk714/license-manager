let connectedServers = [];

console.log('JavaScript file loaded successfully');

// Test function to verify JavaScript is working
function testFunction() {
    console.log('Test function called successfully!');
    alert('JavaScript is working!');
}

function getServerConfig() {
    const hostPort = document.getElementById('server_ip').value;
    const [host, port] = hostPort.includes(':') ? hostPort.split(':') : [hostPort, '22'];
    
    return {
        host: host,
        port: port,
        username: document.getElementById('username').value,
        password: document.getElementById('password').value
    };
}

function generateServerId(host, port) {
    return `${host}_${port}`.replace(/[^a-zA-Z0-9_]/g, '_');
}

function addServer() {
    const config = getServerConfig();
    
    // Validate inputs
    if (!config.host || !config.port || !config.username || !config.password) {
        showStatus('Please fill in all server connection fields', 'error');
        return;
    }

    const serverId = generateServerId(config.host, config.port);
    
    // Check if server already exists
    if (connectedServers.find(s => s.id === serverId)) {
        showStatus('Server already exists in the list', 'error');
        return;
    }

    // Add server to list
    const server = {
        id: serverId,
        host: config.host,
        port: config.port,
        username: config.username,
        password: config.password,
        status: 'checking',
        connected: false
    };
    
    connectedServers.push(server);
    renderServerCard(server);
    updateBatchOperations();
    
    // Clear form
    document.getElementById('server_ip').value = '';
    document.getElementById('password').value = '';
    
    // Check license2_cli on the server
    checkServerLicenseCLI(server);
}

function renderServerCard(server) {
    const container = document.getElementById('connected_servers');
    const cardId = `server_${server.id}`;
    
    // Remove existing card if it exists
    const existingCard = document.getElementById(cardId);
    if (existingCard) {
        existingCard.remove();
    }
    
    const card = document.createElement('div');
    card.id = cardId;
    card.className = `server-card ${server.connected ? 'connected' : server.status === 'error' ? 'error' : ''}`;
    
                card.innerHTML = 
                '<div class="server-card-header">' +
                    '<div class="server-info">' + server.host + ':' + server.port + '</div>' +
                    '<div class="server-status ' + server.status + '">' + server.status + '</div>' +
                '</div>' +
                '<div class="server-actions">' +
                    '<button class="btn btn-sm" onclick="removeServer(\'' + server.id + '\')">' +
                        'Remove' +
                    '</button>' +
                '</div>';
    
    container.appendChild(card);
}

function removeServer(serverId) {
    connectedServers = connectedServers.filter(s => s.id !== serverId);
    const card = document.getElementById(`server_${serverId}`);
    if (card) {
        card.remove();
    }
    updateBatchOperations();
}

async function checkServerLicenseCLI(server) {
    const card = document.getElementById(`server_${server.id}`);
    if (!card) return;
    
    try {
        const response = await fetch('/license-manager/api/check-license-cli', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                host: server.host,
                port: server.port,
                username: server.username,
                password: server.password
            })
        });

        const result = await response.json();

        if (result.exists) {
            server.status = 'connected';
            server.connected = true;
            showStatus(`✓ license2_cli found on ${server.host}:${server.port}`, 'success');
        } else {
            server.status = 'error';
            server.connected = false;
            showStatus(`✗ license2_cli not found on ${server.host}:${server.port}`, 'error');
        }
    } catch (error) {
        server.status = 'error';
        server.connected = false;
        showStatus(`Error checking ${server.host}:${server.port}: ${error.message}`, 'error');
    }
    
    renderServerCard(server);
    updateBatchOperations();
}

function updateBatchOperations() {
    const connectedCount = connectedServers.filter(s => s.connected).length;
    const summary = document.getElementById('batch_server_summary');
    const batchSection = document.getElementById('batch_operations');
    
    if (connectedCount > 0) {
        summary.innerHTML = 
            '<p><strong>' + connectedCount + '</strong> connected server(s) ready for batch operations:</p>' +
            '<ul>' +
                connectedServers.filter(s => s.connected).map(s => 
                    '<li>' + s.host + ':' + s.port + '</li>'
                ).join('') +
            '</ul>';
        batchSection.style.display = 'block';
    } else {
        summary.innerHTML = '<p>No servers connected. Add servers in the Server Connections section above.</p>';
        batchSection.style.display = 'none';
    }
}

async function downloadFromServer(serverId) {
    const server = connectedServers.find(s => s.id === serverId);
    if (!server || !server.connected) return;
    
    try {
        const response = await fetch('/license-manager/api/download-sysinfo', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                host: server.host,
                port: server.port,
                username: server.username,
                password: server.password
            })
        });

        if (response.ok) {
            const contentType = response.headers.get('Content-Type');
            const contentDisposition = response.headers.get('Content-Disposition');
            
            if (contentDisposition && contentType === 'application/octet-stream') {
                const filenameMatch = contentDisposition.match(/filename="(.+)"/);
                const filename = filenameMatch ? filenameMatch[1] : 'sysinfo_file';
                
                const blob = await response.blob();
                const url = window.URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.href = url;
                a.download = filename;
                document.body.appendChild(a);
                a.click();
                window.URL.revokeObjectURL(url);
                document.body.removeChild(a);
                
                showStatus(`✓ Downloaded ${filename} from ${server.host}:${server.port}`, 'success');
            } else {
                const result = await response.json();
                showStatus(`✗ Download failed from ${server.host}:${server.port}: ${result.error}`, 'error');
            }
        } else {
            const result = await response.json();
            showStatus(`✗ Download failed from ${server.host}:${server.port}: ${result.error}`, 'error');
        }
    } catch (error) {
        showStatus(`Error downloading from ${server.host}:${server.port}: ${error.message}`, 'error');
    }
}

function showStatus(message, type = 'info') {
    const statusDiv = document.getElementById('status');
    // Convert line breaks to HTML
    const htmlMessage = message.replace(/\n/g, '<br>');
    statusDiv.innerHTML = htmlMessage;
    statusDiv.className = `status ${type}`;
    statusDiv.classList.remove('hidden');
}

function hideStatus() {
    document.getElementById('status').classList.add('hidden');
}

function setLoading(element, loadingId, isLoading) {
    if (isLoading) {
        element.disabled = true;
        element.style.opacity = '0.6';
    } else {
        element.disabled = false;
        element.style.opacity = '1';
    }
}

function updateFileButton() {
    const fileInput = document.getElementById('license_file');
    const buttonText = document.getElementById('file_button_text');
    const fileWrapper = document.querySelector('.file-input-wrapper');
    
    if (fileInput.files.length > 0) {
        const fileName = fileInput.files[0].name;
        buttonText.textContent = `Upload: ${fileName}`;
        fileWrapper.style.background = 'linear-gradient(135deg, #27ae60 0%, #2ecc71 100%)';
        
        // Auto-upload the file
        uploadLicense();
    } else {
        buttonText.textContent = 'Choose License File';
        fileWrapper.style.background = 'linear-gradient(135deg, #f39c12 0%, #e67e22 100%)';
    }
}

let uploadLineCounter = 0;

function addUploadLine() {
    const uploadLines = document.getElementById('upload_lines');
    const lineId = `upload_line_${uploadLineCounter++}`;
    
    const line = document.createElement('div');
    line.id = lineId;
    line.className = 'upload-line';
    
    // Create server dropdown
    const serverDropdown = document.createElement('select');
    serverDropdown.className = 'server-dropdown';
    serverDropdown.id = `server_${lineId}`;
    
    // Add connected servers to dropdown
    const connectedServers = getConnectedServers();
    if (connectedServers.length === 0) {
        const option = document.createElement('option');
        option.value = '';
        option.textContent = 'No connected servers';
        option.disabled = true;
        serverDropdown.appendChild(option);
    } else {
        connectedServers.forEach(server => {
            const option = document.createElement('option');
            option.value = server.id;
            option.textContent = `${server.host}:${server.port}`;
            serverDropdown.appendChild(option);
        });
    }
    
    // Create file input
    const fileInputDiv = document.createElement('div');
    fileInputDiv.className = 'file-input-small';
    
    const fileInput = document.createElement('input');
    fileInput.type = 'file';
    fileInput.accept = '.lic,.license,.txt';
    fileInput.id = `file_${lineId}`;
    fileInput.onchange = function() { updateFileDisplay(lineId); };
    
    const fileDisplay = document.createElement('div');
    fileDisplay.className = 'file-display';
    fileDisplay.id = `display_${lineId}`;
    fileDisplay.textContent = 'Choose License File';
    fileDisplay.onclick = function() { fileInput.click(); };
    
    fileInputDiv.appendChild(fileInput);
    fileInputDiv.appendChild(fileDisplay);
    
    // Create remove button
    const removeBtn = document.createElement('button');
    removeBtn.type = 'button';
    removeBtn.className = 'remove-line-btn';
    removeBtn.textContent = 'Remove';
    removeBtn.onclick = function() { removeUploadLine(lineId); };
    
    line.appendChild(serverDropdown);
    line.appendChild(fileInputDiv);
    line.appendChild(removeBtn);
    
    uploadLines.appendChild(line);
}

function removeUploadLine(lineId) {
    const line = document.getElementById(lineId);
    if (line) {
        line.remove();
    }
}

function updateFileDisplay(lineId) {
    const fileInput = document.getElementById(`file_${lineId}`);
    const fileDisplay = document.getElementById(`display_${lineId}`);
    
    if (fileInput.files.length > 0) {
        const fileName = fileInput.files[0].name;
        fileDisplay.textContent = fileName;
        fileDisplay.classList.add('has-file');
    } else {
        fileDisplay.textContent = 'Choose License File';
        fileDisplay.classList.remove('has-file');
    }
}

function getConnectedServers() {
    return connectedServers.filter(s => s.connected);
}

function showBatchDownloadStatus(message, type = 'info') {
    const statusDiv = document.getElementById('batch_download_status');
    // Convert line breaks to HTML
    const htmlMessage = message.replace(/\n/g, '<br>');
    statusDiv.innerHTML = htmlMessage;
    statusDiv.className = `status ${type}`;
    statusDiv.classList.remove('hidden');
}

function showBatchUploadStatus(message, type = 'info') {
    console.log('showBatchUploadStatus called with:', message, type);
    console.log('Message length:', message.length);
    console.log('Message contains ```:', message.includes('```'));
    console.log('Message preview:', message.substring(0, 200) + '...');
    
    const statusDiv = document.getElementById('batch_upload_status');
    console.log('statusDiv found:', statusDiv);
    
    // Check if message contains code blocks
    if (message.includes('```')) {
        console.log('Message contains code blocks, converting...');
        // Convert markdown-style code blocks to HTML
        const htmlMessage = message
            .replace(/```\n?([\s\S]*?)\n?```/g, '<pre>$1</pre>')
            .replace(/\n/g, '<br>');
        console.log('Converted HTML:', htmlMessage);
        statusDiv.innerHTML = htmlMessage;
    } else {
        console.log('Message does not contain code blocks');
        // Convert line breaks to HTML
        const htmlMessage = message.replace(/\n/g, '<br>');
        statusDiv.innerHTML = htmlMessage;
    }
    statusDiv.className = `status ${type}`;
    statusDiv.classList.remove('hidden');
    console.log('Status div updated, class:', statusDiv.className);
}

function hideBatchDownloadStatus() {
    document.getElementById('batch_download_status').classList.add('hidden');
}

function hideBatchUploadStatus() {
    document.getElementById('batch_upload_status').classList.add('hidden');
}

async function batchDownloadSysinfo() {
    const servers = getConnectedServers();
    if (servers.length === 0) {
        showBatchDownloadStatus('No connected servers available for batch download', 'error');
        return;
    }

    const button = document.querySelector('button[onclick="batchDownloadSysinfo()"]');
    if (!button) {
        console.error('Batch download button not found');
        return;
    }
    setLoading(button, 'batch_download_loading', true);
    hideBatchDownloadStatus();

    let successCount = 0;
    let errorCount = 0;
    const results = [];

    for (let i = 0; i < servers.length; i++) {
        const server = servers[i];
        try {
            const response = await fetch('/license-manager/api/download-sysinfo', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    host: server.host,
                    port: server.port,
                    username: server.username,
                    password: server.password
                })
            });

            if (response.ok) {
                const contentType = response.headers.get('Content-Type');
                const contentDisposition = response.headers.get('Content-Disposition');
                
                if (contentDisposition && contentType === 'application/octet-stream') {
                    const filenameMatch = contentDisposition.match(/filename="(.+)"/);
                    const filename = filenameMatch ? filenameMatch[1] : 'sysinfo_file';
                    
                    const blob = await response.blob();
                    const url = window.URL.createObjectURL(blob);
                    const a = document.createElement('a');
                    a.href = url;
                    a.download = filename;
                    document.body.appendChild(a);
                    a.click();
                    window.URL.revokeObjectURL(url);
                    document.body.removeChild(a);
                    
                    results.push(`✓ ${server.host}:${server.port} - ${filename}`);
                    successCount++;
                } else {
                    const result = await response.json();
                    results.push(`✗ ${server.host}:${server.port} - ${result.error}`);
                    errorCount++;
                }
            } else {
                const result = await response.json();
                results.push(`✗ ${server.host}:${server.port} - ${result.error}`);
                errorCount++;
            }
        } catch (error) {
            results.push(`✗ ${server.host}:${server.port} - ${error.message}`);
            errorCount++;
        }
    }

    const summary = `Batch download completed: ${successCount} successful, ${errorCount} failed\n\n${results.join('\n')}`;
    showBatchDownloadStatus(summary, errorCount === 0 ? 'success' : 'error');
    setLoading(button, 'batch_download_loading', false);
}

async function uploadAllFiles() {
    console.log('uploadAllFiles called');
    const uploadLines = document.querySelectorAll('.upload-line');
    console.log('Found upload lines:', uploadLines.length);
    
    if (uploadLines.length === 0) {
        console.log('No upload lines found, showing error');
        showBatchUploadStatus('Please add at least one file to upload', 'error');
        return;
    }

    // Validate all upload lines
    const uploadTasks = [];
    for (let i = 0; i < uploadLines.length; i++) {
        const line = uploadLines[i];
        const lineId = line.id;
        
        const serverSelect = document.getElementById(`server_${lineId}`);
        const fileInput = document.getElementById(`file_${lineId}`);
        
        if (!serverSelect.value) {
            showBatchUploadStatus('Please select a server for all upload lines', 'error');
            return;
        }
        
        if (!fileInput.files[0]) {
            showBatchUploadStatus('Please select a file for all upload lines', 'error');
            return;
        }
        
        // Find the server details
        const server = connectedServers.find(s => s.id === serverSelect.value);
        if (!server) {
            showBatchUploadStatus('Selected server not found', 'error');
            return;
        }
        
        uploadTasks.push({
            server: server,
            file: fileInput.files[0],
            lineId: lineId
        });
    }

    const button = document.querySelector('button[onclick="uploadAllFiles()"]');
    console.log('Upload button found:', button);
    if (!button) {
        console.error('Upload button not found');
        return;
    }
    console.log('Setting loading state and hiding status');
    setLoading(button, 'upload_loading', true);
    hideBatchUploadStatus();

    // Show upload progress list
    const progressList = document.getElementById('upload_progress_list');
    const uploadItems = document.getElementById('upload_items');
    progressList.classList.remove('hidden');
    uploadItems.innerHTML = '';

    // Create upload items for each task
    uploadTasks.forEach((task, index) => {
        const item = document.createElement('div');
        item.id = `upload_${task.lineId}`;
        item.className = 'upload-item processing';
        item.innerHTML = 
            '<span>' + task.server.host + ':' + task.server.port + ' - ' + task.file.name + '</span>' +
            '<span>Processing...</span>';
        uploadItems.appendChild(item);
    });

    let successCount = 0;
    let errorCount = 0;
    const results = [];
    const detailedMessages = [];

    for (let i = 0; i < uploadTasks.length; i++) {
        const task = uploadTasks[i];
        const item = document.getElementById(`upload_${task.lineId}`);
        
        try {
            const formData = new FormData();
            formData.append('license_file', task.file);
            formData.append('host', task.server.host);
            formData.append('port', task.server.port);
            formData.append('username', task.server.username);
            formData.append('password', task.server.password);

            const response = await fetch('/license-manager/api/upload-license', {
                method: 'POST',
                body: formData
            });

            const result = await response.json();

            if (result.success) {
                item.className = 'upload-item success';
                item.innerHTML = 
                    '<span>' + task.server.host + ':' + task.server.port + ' - ' + task.file.name + '</span>' +
                    '<span>✓ Success</span>';
                results.push(`✓ ${task.server.host}:${task.server.port} - ${task.file.name} uploaded successfully`);
                if (result.message) {
                    detailedMessages.push(`\n--- ${task.server.host}:${task.server.port} - ${task.file.name} ---\n${result.message}`);
                }
                successCount++;
            } else {
                item.className = 'upload-item error';
                item.innerHTML = 
                    '<span>' + task.server.host + ':' + task.server.port + ' - ' + task.file.name + '</span>' +
                    '<span>✗ ' + result.error + '</span>';
                results.push(`✗ ${task.server.host}:${task.server.port} - ${task.file.name} - ${result.error}`);
                errorCount++;
            }
        } catch (error) {
            item.className = 'upload-item error';
            item.innerHTML = 
                '<span>' + task.server.host + ':' + task.server.port + ' - ' + task.file.name + '</span>' +
                '<span>✗ ' + error.message + '</span>';
            results.push(`✗ ${task.server.host}:${task.server.port} - ${task.file.name} - ${error.message}`);
            errorCount++;
        }
    }

    let summary = `Upload completed: ${successCount} successful, ${errorCount} failed\n\n${results.join('\n')}`;
    if (detailedMessages.length > 0) {
        summary += '\n\n' + detailedMessages.join('\n');
    }
    console.log('Final summary:', summary);
    console.log('Calling showBatchUploadStatus with summary length:', summary.length);
    showBatchUploadStatus(summary, errorCount === 0 ? 'success' : 'error');
    setLoading(button, 'upload_loading', false);
}
