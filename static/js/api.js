function registerApiComponents() {
    if (window.ApiComponentsRegistered) return;
    window.ApiComponentsRegistered = true;

    Alpine.data('apiUserGet', () => ({
        user: null, 
        loading: false,
        async fetchUser() {
            const id = this.$el.querySelector('input[name=user_id]').value;
            if(!id) return;
            this.loading = true;
            this.user = null;
            try {
                const response = await window.apiFetch('/api/users/' + id, { 
                    credentials: 'include',
                    headers: window.getAuthHeaders()
                });
                if (response.ok) {
                    this.user = await response.json();
                } else {
                    let errorMsg = 'Пользователь не найден';
                    try {
                        const errorData = await response.json();
                        errorMsg = errorData.error || errorMsg;
                    } catch(e) {}
                    this.$dispatch('notify', { message: errorMsg, type: 'error' });
                }
            } catch (e) {
                this.$dispatch('notify', { message: 'Ошибка сети', type: 'error' });
            } finally {
                this.loading = false;
            }
        }
    }));

    Alpine.data('apiUserList', () => ({
        users: [], 
        loading: false,
        async fetchUsers() {
            this.loading = true;
            this.users = [];
            try {
                const response = await window.apiFetch('/api/users', { 
                    credentials: 'include',
                    headers: window.getAuthHeaders()
                });
                if (response.ok) {
                    this.users = await response.json();
                } else {
                    let errorMsg = 'Доступ запрещен или ошибка сервера';
                    try {
                        const errorData = await response.json();
                        errorMsg = errorData.error || errorMsg;
                    } catch(e) {}
                    this.$dispatch('notify', { message: errorMsg, type: 'error' });
                }
            } catch (e) {
                this.$dispatch('notify', { message: 'Ошибка загрузки', type: 'error' });
            } finally {
                this.loading = false;
            }
        }
    }));
}

document.addEventListener('alpine:init', registerApiComponents);
if (window.Alpine) registerApiComponents();
