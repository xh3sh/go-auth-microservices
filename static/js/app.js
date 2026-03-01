window.getAuthHeaders = function() {
    const headers = {};
    
    const basicUser = document.getElementById('basic-username')?.value;
    const basicPass = document.getElementById('basic-password')?.value;

    if (basicUser && basicPass) {
        const auth = btoa(`${basicUser}:${basicPass}`);
        headers['Authorization'] = `Basic ${auth}`;
    } else {
        const token = sessionStorage.getItem('basic_token');
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }
    }
    return headers;
};

// apiFetch выполняет сетевой запрос с автоматическим обновлением JWT токена
window.apiFetch = async function(url, options = {}) {
    let response = await fetch(url, options);

    if (response.status === 401) {
        const isBasic = options.headers && options.headers['Authorization'] && options.headers['Authorization'].startsWith('Basic');
        const hasJwtHint = sessionStorage.getItem('basic_token') || document.cookie.includes('jwt_logged_in=true');

        if (!isBasic && hasJwtHint) {
            try {
                const refreshResponse = await fetch('/auth/jwt/refresh', { 
                    method: 'POST', 
                    credentials: 'include' 
                });

                if (refreshResponse.ok) {
                    const data = await refreshResponse.json();
                    
                    if (sessionStorage.getItem('basic_token')) {
                        sessionStorage.setItem('basic_token', data.access_token);
                        
                        if (options.headers) {
                            options.headers['Authorization'] = `Bearer ${data.access_token}`;
                        }
                    }

                    return await fetch(url, options);
                } else {
                    if (refreshResponse.status === 401 || refreshResponse.status === 400) {
                        document.cookie = "jwt_logged_in=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
                    }
                }
            } catch (e) {
                console.error('Ошибка при обновлении токена:', e);
            }
        }
    }

    return response;
};

function registerGlobalComponents() {
    if (window.GlobalComponentsRegistered) return;
    window.GlobalComponentsRegistered = true;

    Alpine.data('notificationSystem', () => ({
        notifications: [],
        addNotification(message, type = 'success') {
            const id = Date.now();
            this.notifications.push({ id, message, type });
            setTimeout(() => {
                this.notifications = this.notifications.filter(n => n.id !== id);
            }, 5000);
        }
    }));
}

document.addEventListener('alpine:init', registerGlobalComponents);
if (window.Alpine) registerGlobalComponents();

document.addEventListener('DOMContentLoaded', () => {
    document.body.addEventListener('htmx:configRequest', (event) => {
        event.detail.withCredentials = true;
        const authHeaders = window.getAuthHeaders();
        Object.assign(event.detail.headers, authHeaders);
    });

    document.body.addEventListener('htmx:afterSettle', (event) => {
        const target = event.detail.target;
        if (target.id === 'auth-result' || target.id === 'log-result') {
            try {
                const text = target.textContent.trim();
                if (!text || text.startsWith('<')) return;
                const json = JSON.parse(text);
                
                let fieldsHtml = '';
                for (const [key, value] of Object.entries(json)) {
                    if (typeof value === 'object' || key === 'access_token' || key === 'token') continue;
                    fieldsHtml += `<div class="ui-field"><span class="ui-field-label">${key}</span><span class="ui-field-value">${value}</span></div>`;
                }

                const isSuccess = !json.error && (json.valid || json.authenticated || json.username || json.user_id);

                target.innerHTML = `
                    <div class="ui-card ${isSuccess ? 'card-user' : 'card-auth-success'}" style="margin-top: 20px; max-width: 100%;">
                        <div class="ui-card-header">
                            <span class="ui-card-title">${json.username || json.user || json.message || 'Результат'}</span>
                            <span class="ui-card-badge">${isSuccess ? 'Success' : 'Info'}</span>
                        </div>
                        <div class="ui-card-body">${fieldsHtml}</div>
                    </div>`;
            } catch (e) {}
        }
    });
});
