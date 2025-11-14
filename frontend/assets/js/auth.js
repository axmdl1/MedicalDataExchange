// Authentication Manager
class AuthManager {
    constructor() {
        this.token = localStorage.getItem('auth_token');
        this.user = this.loadUser();
    }

    loadUser() {
        const userData = localStorage.getItem('user_data');
        return userData ? JSON.parse(userData) : null;
    }

    saveAuth(token, user) {
        this.token = token;
        this.user = user;
        localStorage.setItem('auth_token', token);
        localStorage.setItem('user_data', JSON.stringify(user));
    }

    clearAuth() {
        this.token = null;
        this.user = null;
        localStorage.removeItem('auth_token');
        localStorage.removeItem('user_data');
    }

    isAuthenticated() {
        return !!this.token;
    }

    getToken() {
        return this.token;
    }

    getUser() {
        return this.user;
    }

    getRole() {
        return this.user ? this.user.type : null;
    }

    getClinicId() {
        // clinic_id может быть объектом с полем value (protobuf Int64Value)
        if (!this.user) return null;
        if (this.user.clinic_id && typeof this.user.clinic_id === 'object' && 'value' in this.user.clinic_id) {
            return this.user.clinic_id.value;
        }
        return this.user.clinic_id;
    }

    isPatient() {
        return this.getRole() === ROLES.PATIENT;
    }

    isEmployee() {
        return this.getRole() === ROLES.EMPLOYEE;
    }

    isAdmin() {
        return this.getRole() === ROLES.ADMIN;
    }
}

// Global auth instance
const auth = new AuthManager();
