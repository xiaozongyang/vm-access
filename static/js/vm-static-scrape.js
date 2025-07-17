function validateVmStaticScrape(data) {
    const targets = data.endpoint.targets;

    const nameRegex = /^[a-zA-Z_][a-zA-Z0-9_-]*$/;
    if (!nameRegex.test(data.name)) {
        alert(`Name "${data.name}" is invalid, must be in a-zA-Z_- format`);
        return false;
    }

    const targetRegex = /^[^:]+:\d+$/;
    const httpRegex = /^https?:\/\//i;

    for (const target of targets) {
        if (httpRegex.test(target)) {
            alert(`Target "${target}" can not contain http:// or https:// protocol`);
            targetsInput.focus();
            return false;
        }

        if (!targetRegex.test(target)) {
            alert(`Target "${target}" is invalid, must be in ip:port format (like 1.2.3.4:9100)`);
            targetsInput.focus();
            return false;
        }
    }

    const labels = data.endpoint.labels;

    const labelRegex = /^[a-zA-Z_][a-zA-Z0-9_]*$/;

    for (const [key, value] of Object.entries(labels)) {
        if (key === undefined || value === undefined) {
            alert(`Label "${key}=${value}" is invalid, must be in key=value format`);
            return false;
        }
        if (!labelRegex.test(key)) {
            alert(`Label key "${key}" is invalid, must be in a-zA-Z_ format`);
            return false;
        }
        if (!labelRegex.test(value)) {
            alert(`Label value "${value}" is invalid, must be in a-zA-Z_ format`);
            return false;
        }
    }

    return true;
}

function createVmStaticScrape(url) {
    return createOrUpdateVmStaticScrape('POST', url);
}

function updateVmStaticScrape(url) {
    return createOrUpdateVmStaticScrape('PUT', url);
}

function createOrUpdateVmStaticScrape(method, url) {
    const data = readVmStaticScrapeFormData();

    if (!validateVmStaticScrape(data)) {
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
            alert('create or update VMStaticScrape successfully!');
            window.location.href = getStaticScrapeRedirectPath(method, data.name);
        } else {
            return response.json().then(data => {
                throw new Error(data.error || 'Failed to create or update VM Static Scrape');
            });
        }
    }).catch(error => {
        alert('Error: ' + error.message);
    });
}

function deleteVmStaticScrape(url) {
    const data = readVmStaticScrapeFormData();

    const isConfirmed = confirm('VMStaticScrape ' + data.name + ' will be deleted, this action cannot be undone!');
    if (!isConfirmed) {
        return;
    }

    fetch(url, {
        method: 'DELETE',
    })
        .then(response => {
            if (response.ok) {
                alert('VM Static Scrape deleted successfully!');
                window.location.href = getStaticScrapeRedirectPath('DELETE', data.name);
            }
        });
}

function readVmStaticScrapeFormData() {
    return {
        name: document.getElementById('name').value,
        cluster: getCluster(),
        jobName: document.getElementById('jobName').value,
        endpoint: {
            path: document.getElementById('path').value || '/metrics',
            targets: document.getElementById('targets').value.split('\n')
                .map(line => line.trim())
                .filter(line => line !== ''),
            labels: document.getElementById('labels').value
                .split('\n')
                .filter(line => line.includes('='))
                .reduce((acc, line) => {
                    const [key, value] = line.split('=');
                    if (key && value) {
                        acc[key.trim()] = value.trim();
                    }
                    return acc;
                }, {}),
        },
        meta: {
            owner: document.getElementById('owner').value,
        },
    };
}

function getStaticScrapeRedirectPath(method, name) {
    return getRedirectPath(method, 'vm-static-scrapes', name);
}