function getCluster() {
    return JSON.parse(document.getElementById("page-data").textContent).cluster;
}

function getRedirectPath(method, resource, id) {
    let clusterPath = '/clusters/' + getCluster();
    if (method === 'DELETE') {
        return clusterPath;
    }
    return clusterPath + '/' + resource + '/' + id;
}

function toggleAdvancedOptions() {
    const showAdvanced = document.getElementById('showAdvanced');
    const advancedOptions = document.getElementById('advancedOptions');
    advancedOptions.classList.toggle('hidden', !showAdvanced.checked);

    showAdvanced.addEventListener('change', function () {
        advancedOptions.classList.toggle('hidden', !showAdvanced.checked);
    });
}