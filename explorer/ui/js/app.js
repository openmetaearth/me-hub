// ME Network Explorer - Real-time Block Height Display
// Polls the backend API every 5 seconds for latest block height

class ExplorerApp {
    constructor() {
        this.apiBaseUrl = 'http://localhost:8080';
        this.pollInterval = 5000; // 5 seconds
        this.blockHistory = [];
        this.maxHistoryPoints = 20;
        this.chart = null;
        this.previousHeight = null;
        
        this.init();
    }

    async init() {
        this.setupChart();
        await this.fetchData();
        this.startPolling();
    }

    setupChart() {
        const ctx = document.getElementById('blockChart').getContext('2d');
        this.chart = new Chart(ctx, {
            type: 'line',
            data: {
                labels: [],
                datasets: [{
                    label: 'Block Height',
                    data: [],
                    borderColor: '#38bdf8',
                    backgroundColor: 'rgba(56, 189, 248, 0.1)',
                    fill: true,
                    tension: 0.4,
                    pointRadius: 3,
                    pointBackgroundColor: '#38bdf8'
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                animation: {
                    duration: 300
                },
                scales: {
                    x: {
                        display: true,
                        grid: {
                            color: 'rgba(148, 163, 184, 0.1)'
                        },
                        ticks: {
                            color: '#94a3b8',
                            maxTicksLimit: 10
                        }
                    },
                    y: {
                        display: true,
                        grid: {
                            color: 'rgba(148, 163, 184, 0.1)'
                        },
                        ticks: {
                            color: '#94a3b8',
                            callback: function(value) {
                                return value.toLocaleString();
                            }
                        }
                    }
                },
                plugins: {
                    legend: {
                        display: false
                    }
                }
            }
        });
    }

    async fetchData() {
        try {
            const [blockHeightResponse, statusResponse] = await Promise.all([
                fetch(`${this.apiBaseUrl}/api/block-height`),
                fetch(`${this.apiBaseUrl}/api/status`)
            ]);

            if (!blockHeightResponse.ok || !statusResponse.ok) {
                throw new Error('Failed to fetch data from API');
            }

            const blockData = await blockHeightResponse.json();
            const statusData = await statusResponse.json();

            this.updateUI(blockData, statusData);
        } catch (error) {
            console.error('Error fetching data:', error);
            this.showError();
        }
    }

    updateUI(blockData, statusData) {
        const height = blockData.height;
        const timestamp = new Date(blockData.timestamp);
        
        // Update block height display with animation
        const heightElement = document.getElementById('blockHeight');
        heightElement.textContent = height.toLocaleString();
        heightElement.classList.add('updating');
        setTimeout(() => heightElement.classList.remove('updating'), 500);

        // Update latest block card
        document.getElementById('latestHeight').textContent = height.toLocaleString();
        document.getElementById('latestTime').textContent = timestamp.toLocaleTimeString();

        // Update sync status
        const syncInfo = statusData.sync_info;
        const syncStatusElement = document.getElementById('syncStatus');
        if (syncInfo.catching_up) {
            syncStatusElement.textContent = 'Catching Up';
            syncStatusElement.className = 'sync-status catching-up';
        } else {
            syncStatusElement.textContent = 'Synced';
            syncStatusElement.className = 'sync-status synced';
        }

        // Update network status
        document.getElementById('networkStatus').textContent = syncInfo.catching_up ? 'Catching Up' : 'Running';
        document.getElementById('networkStatus').className = syncInfo.catching_up ? 'value catching-up' : 'value synced';
        document.getElementById('catchingUp').textContent = syncInfo.catching_up ? 'Yes' : 'No';

        // Calculate blocks per minute
        if (this.previousHeight !== null) {
            const heightDiff = height - this.previousHeight;
            const blocksPerMin = (heightDiff / (this.pollInterval / 60000)).toFixed(1);
            document.getElementById('blocksPerMin').textContent = blocksPerMin;
        }
        this.previousHeight = height;

        // Update last update time
        document.getElementById('lastUpdate').textContent = new Date().toLocaleTimeString();

        // Update chart
        this.updateChart(height, timestamp);
    }

    updateChart(height, timestamp) {
        const time = timestamp.toLocaleTimeString();
        
        this.blockHistory.push({
            height: height,
            time: time
        });

        if (this.blockHistory.length > this.maxHistoryPoints) {
            this.blockHistory.shift();
        }

        this.chart.data.labels = this.blockHistory.map(d => d.time);
        this.chart.data.datasets[0].data = this.blockHistory.map(d => d.height);
        this.chart.update('none');
    }

    showError() {
        document.getElementById('blockHeight').textContent = 'Error';
        document.getElementById('syncStatus').textContent = 'Connection Error';
        document.getElementById('syncStatus').className = 'sync-status error';
        document.getElementById('networkStatus').textContent = 'Unknown';
        document.getElementById('networkStatus').className = 'value error';
        document.getElementById('catchingUp').textContent = 'Unknown';
        document.getElementById('latestHeight').textContent = 'Error';
        document.getElementById('latestTime').textContent = 'Error';
        document.getElementById('blocksPerMin').textContent = '--';
    }

    startPolling() {
        setInterval(() => {
            this.fetchData();
        }, this.pollInterval);
    }
}

// Initialize the app when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    new ExplorerApp();
});
