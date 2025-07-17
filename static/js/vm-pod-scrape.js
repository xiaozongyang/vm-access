
function createVMPodScrape(url) {
    return createOrUpdateVMPodScrape('POST', url);
}

function updateVMPodScrape(url) {
    return createOrUpdateVMPodScrape('PUT', url);
}

function createOrUpdateVMPodScrape(method, url) {
    const data = readVMPodScrapeFormData();
    if (!validateVMPodScrape(data)) {
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
            alert('create or update VMPodScrape successfully!');
            window.location.href = getPodScrapeRedirectPath(method, data.name);
        } else {
            return response.json().then(data => {
                throw new Error(data.error || 'Failed to create or update VMPodScrape');
            });
        }
    }).catch(error => {
        alert(`Failed to create or update VMPodScrape: ${error}`);
    });
}

function deleteVMPodScrape(url) {
    fetch(url, {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json',
        },
    }).then(response => {
        if (response.ok) {
            alert('delete VMPodScrape successfully!');
            window.location.href = getPodScrapeRedirectPath('DELETE', data.name);
        } else {
            return response.json().then(data => {
                throw new Error(data.error || 'Failed to delete VMPodScrape');
            });
        }
    }).catch(error => {
        alert(`Failed to delete VMPodScrape: ${error}`);
    });
}

function readVMPodScrapeFormData() {
    const selectorType = document.getElementById('selectorType').value;

    const data = {
        name: document.getElementById('name').value,
        cluster: getCluster(),
        namespace: document.getElementById('namespace').value,
        selector: {},
        port: document.getElementById('port').value,
        path: document.getElementById('path').value,
        jobLabel: document.getElementById('jobLabel').value.trim(),
        meta: {
            owner: document.getElementById('owner').value,
        },
    };

    if (selectorType === 'matchLabels') {
        data.selector.matchLabels = document.getElementById('selectorMatchLabels').value
            .split('\n')
            .map(line => line.trim())
            .filter(line => line.includes('='))
            .reduce((acc, line) => {
                const [key, value] = line.split('=');
                if (key && value) {
                    acc[key.trim()] = value.trim();
                }
                return acc;
            }, {});
    } else if (selectorType === 'matchExpressions') {
        data.selector.matchExpressions = document.getElementById('selectorMatchExpressions').value
            .split('\n')
            .map(line => line.trim())
            .filter(line => line.includes(' '))
            .map(line => {
                const [key, operator, values] = line.split(' ');
                return {
                    key: key.trim(),
                    operator: operator.trim(),
                    values: values.trim().replace(/[\[\]]/g, '').split(',').map(value => value.trim()),
                };
            });
    }

    return data;
}

function getPodScrapeRedirectPath(method, name) {
    return getRedirectPath(method, 'vm-pod-scrapes', name);
}

function validateVMPodScrape(data) {
    const nameRegex = /^[a-zA-Z_][a-zA-Z0-9_-]*$/;
    const selectorKeyRegex = /^(?:(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\/)?[a-zA-Z0-9](?:[a-zA-Z0-9-_.]{0,61}[a-zA-Z0-9])?$/;
    const selectorValueRegex = /^[a-zA-Z0-9](?:[a-zA-Z0-9-_.]{0,61}[a-zA-Z0-9])?$/;

    if (!nameRegex.test(data.name)) {
        alert(`Name "${data.name}" is invalid, must be in a-zA-Z_- format`);
        return false;
    }

    if (!data.selector || (Object.keys(data.selector).length === 0)) {
        alert('Selector must be provided');
        return false;
    }

    if (data.selector.matchLabels) {
        for (const [key, value] of Object.entries(data.selector.matchLabels)) {
            if (!key || !value) {
                alert(`Label "${key}=${value}" is invalid`);
                return false;
            }
            if (!selectorKeyRegex.test(key)) {
                alert(`Label key "${key}" is invalid`);
                return false;
            }
            if (!selectorValueRegex.test(value)) {
                alert(`Label value "${value}" is invalid for key "${key}"`);
                return false;
            }
        }
    }

    if (data.selector.matchExpressions) {
        for (const expression of data.selector.matchExpressions) {
            if (!expression.key || !expression.operator || !expression.values) {
                alert(`Match expression "${expression.key} ${expression.operator} ${expression.values}" is invalid`);
                return false;
            }
            if (expression.operator !== 'In' && expression.operator !== 'NotIn') {
                alert(`Match expression operator "${expression.operator}" is invalid`);
                return false;
            }
            if (expression.values.length === 0) {
                alert(`Match expression values cannot be empty for key "${expression.key}"`);
                return false;
            }
            if (!selectorKeyRegex.test(expression.key)) {
                alert(`Match expression key "${expression.key}" is invalid`);
                return false;
            }
            for (const value of expression.values) {
                if (!selectorValueRegex.test(value)) {
                    alert(`Match expression value "${value}" is invalid for key "${expression.key}"`);
                    return false;
                }
            }
        }
    }

    return true;
}