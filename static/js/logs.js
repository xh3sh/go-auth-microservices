function registerLogComponents() {
    if (window.LogComponentsRegistered) return;
    window.LogComponentsRegistered = true;

    Alpine.data('logViewer', () => ({
        logs: [],
        loading: false,
        activeLogTab: 'all',
        async fetchLogs(url) {
            this.loading = true;
            try {
                const response = await window.apiFetch(url, {
                    credentials: 'include',
                    headers: window.getAuthHeaders()
                });
                if (response.ok) {
                    const data = await response.json();
                    const rawLogs = data.logs || [];
                    this.logs = rawLogs.map(log => {
                        try {
                            const parsed = JSON.parse(log.data);
                            log.formattedData = JSON.stringify(parsed, null, 2);
                        } catch (e) {
                            log.formattedData = log.data;
                        }
                        const d = new Date(log.timestamp);
                        log.displayDate = d.toLocaleString('ru-RU');
                        return log;
                    });
                } else {
                    this.$dispatch('notify', { message: 'Ошибка получения логов', type: 'error' });
                }
            } catch (e) {
                this.$dispatch('notify', { message: 'Ошибка сети', type: 'error' });
            } finally {
                this.loading = false;
            }
        },
        async applyFilters(event) {
            const form = event ? event.target : this.$el.querySelector('form');
            if (!form) return;
            const formData = new FormData(form);
            const params = new URLSearchParams();
            for (const [key, value] of formData.entries()) {
                if (value) params.append(key, value);
            }
            params.append('page_size', '50');
            await this.fetchLogs('/log/logs/filter?' + params.toString());
        }
    }));
}

document.addEventListener('alpine:init', registerLogComponents);
if (window.Alpine) registerLogComponents();
