
function createVmServiceScrape(url) {
    return createOrUpdateVmServiceScrape('POST', url);
}

function updateVmServiceScrape(url) {
    return createOrUpdateVmServiceScrape('PUT', url);
}

function createOrUpdateVmServiceScrape(method, url) {
    const data = readVmServiceScrapeFormData();
    if (!validateVmServiceScrape(data)) {
        return;
    }

    fetch(url, {
        method: method,
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
    }).then(response => {
        if (response.ok) {
            alert('create or update VMServiceScrape successfully!');
            window.location.href = getServiceScrapeRedirectPath(method, data.name);
        } else {
            return response.json().then(data => {
                throw new Error(data.error || 'Failed to create or update VMServiceScrape');
            });
        }
    }).catch(error => {
        alert(`Failed to create or update VMServiceScrape: ${error}`);
    });
}

function deleteVmServiceScrape(url) {
    fetch(url, {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json',
        },
    }).then(response => {
        if (response.ok) {
            alert('delete VMServiceScrape successfully!');
            window.location.href = getServiceScrapeRedirectPath('DELETE', data.name);
        } else {
            return response.json().then(data => {
                throw new Error(data.error || 'Failed to delete VMServiceScrape');
            });
        }
    }).catch(error => {
        alert(`Failed to delete VMServiceScrape: ${error}`);
    });
}

function readVmServiceScrapeFormData() {
    return {
        name: document.getElementById('name').value,
        cluster: getCluster(),
        namespace: document.getElementById('namespace').value,
        selector: document.getElementById('selector').value
            .split('\n')
            .map(line => line.trim())
            .filter(line => line.includes('='))
            .reduce((acc, line) => {
                const [key, value] = line.split('=');
                if (key && value) {
                    acc[key.trim()] = value.trim();
                }
                return acc;
            }, {}),
        port: document.getElementById('port').value,
        path: document.getElementById('path').value,
        jobLabel: document.getElementById('jobLabel').value.trim(),
        meta: {
            owner: document.getElementById('owner').value,
        },
    };
}

function getServiceScrapeRedirectPath(method, name) {
    return getRedirectPath(method, 'vm-service-scrapes', name);
}

function validateVmServiceScrape(data) {
    const nameRegex = /^[a-zA-Z_][a-zA-Z0-9_-]*$/;
    if (!nameRegex.test(data.name)) {
        alert(`Name "${data.name}" is invalid, must be in a-zA-Z_- format`);
        return false;
    }

    // 更新后的正则表达式
    const selectorKeyRegex = /^(?:(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\/)?[a-zA-Z0-9](?:[a-zA-Z0-9-_.]{0,61}[a-zA-Z0-9])?$/;
    const selectorValueRegex = /^[a-zA-Z0-9](?:[a-zA-Z0-9-_.]{0,61}[a-zA-Z0-9])?$/;

    if (!data.selector || Object.keys(data.selector).length < 1) {
        alert('At least one selector label is required.');
        return false;
    }

    for (const [key, value] of Object.entries(data.selector)) {
        if (key === undefined || value === undefined) {
            alert(`Label "${key}=${value}" is invalid, must be in key=value format`);
            return false;
        }

        // Key 验证（支持带前缀的 key）
        if (!selectorKeyRegex.test(key)) {
            alert(`Label key "${key}" is invalid. Valid examples: "app", "app.kubernetes.io/name"`);
            return false;
        }

        // Value 验证（不能为空）
        if (typeof value !== 'string' || value === '') {
            alert(`Label value cannot be empty for key "${key}"`);
            return false;
        }
        if (!selectorValueRegex.test(value)) {
            alert(`Label value "${value}" is invalid for key "${key}". Only alphanumeric, '-', '_' and '.' are allowed`);
            return false;
        }
    }

    return true;
}