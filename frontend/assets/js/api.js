// API Client
class APIClient {
    constructor() {
        this.baseURL = API_CONFIG.BASE_URL;
    }

    async request(endpoint, options = {}) {
        const url = `${this.baseURL}${endpoint}`;
        const headers = {
            'Content-Type': 'application/json',
            ...options.headers
        };

        if (auth.isAuthenticated()) {
            headers['Authorization'] = `Bearer ${auth.getToken()}`;
        }

        const config = {
            ...options,
            headers
        };

        try {
            const response = await fetch(url, config);
            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.error || 'Request failed');
            }

            return data;
        } catch (error) {
            console.error('API Error:', error);
            throw error;
        }
    }

    // LOGIN
    async login(email, password) {
        return this.request(API_CONFIG.ENDPOINTS.LOGIN, {
            method: 'POST',
            body: JSON.stringify({ email, password })
        });
    }

    // USERS
    async getUsers(filters = {}) {
        const params = new URLSearchParams(filters);
        return this.request(`${API_CONFIG.ENDPOINTS.USERS}?${params}`);
    }

    async getUser(id) {
        return this.request(API_CONFIG.ENDPOINTS.USER(id));
    }

    async createUser(userData) {
        const currentUser = auth.isAuthenticated() ? auth.getUser() : null;
        const userToSend = { ...userData };

        if (currentUser && currentUser.role === 'employee') {
            if (!currentUser.clinic_id) {
                throw new Error('Employee has no associated clinic');
            }
            userToSend.clinic_id = currentUser.clinic_id;
        }

        return this.request(API_CONFIG.ENDPOINTS.USERS, {
            method: 'POST',
            body: JSON.stringify({ user: userToSend })
        });
    }

    async updateUser(id, userData) {
        const userToSend = { ...userData };
        if (userData.clinic_id !== undefined && userData.clinic_id !== null) {
            userToSend.clinic_id = { value: userData.clinic_id };
        }

        return this.request(API_CONFIG.ENDPOINTS.USER(id), {
            method: 'PATCH',
            body: JSON.stringify({ user: userToSend })
        });
    }

    async deleteUser(id) {
        return this.request(API_CONFIG.ENDPOINTS.USER(id), {
            method: 'DELETE'
        });
    }

    // MEDICAL DATA
    async getMedicalData(filters = {}) {
        const params = new URLSearchParams(filters);
        return this.request(`${API_CONFIG.ENDPOINTS.MEDICAL_DATA}?${params}`);
    }

    async getMedicalDataItem(id) {
        return this.request(API_CONFIG.ENDPOINTS.MEDICAL_DATA_ITEM(id));
    }

    async createMedicalData(data) {
        return this.request(API_CONFIG.ENDPOINTS.MEDICAL_DATA, {
            method: 'POST',
            body: JSON.stringify(data)
        });
    }

    async deleteMedicalData(id) {
        return this.request(API_CONFIG.ENDPOINTS.MEDICAL_DATA_ITEM(id), {
            method: 'DELETE'
        });
    }

    // TRANSFERS
    async getTransfers(filters = {}) {
        const params = new URLSearchParams(filters);
        return this.request(`${API_CONFIG.ENDPOINTS.TRANSFERS}?${params}`);
    }

    async getTransfer(id) {
        return this.request(API_CONFIG.ENDPOINTS.TRANSFER(id));
    }

    async createTransfer(data) {
        return this.request(API_CONFIG.ENDPOINTS.TRANSFERS, {
            method: 'POST',
            body: JSON.stringify(data)
        });
    }

    async makeTransferDecision(transferId, confirm) {
        const user = auth.getUser();

        return this.request(API_CONFIG.ENDPOINTS.TRANSFER_DECISION, {
            method: 'POST',
            body: JSON.stringify({
                transfer_id: transferId,
                confirm,
                user_id: user.id
            })
        });
    }

    // CLINICS
    async getClinics() {
        return this.request(API_CONFIG.ENDPOINTS.CLINICS);
    }

    async getClinic(id) {
        return this.request(API_CONFIG.ENDPOINTS.CLINIC(id));
    }

    async createClinic(clinicData) {
        return this.request(API_CONFIG.ENDPOINTS.CLINICS, {
            method: 'POST',
            body: JSON.stringify({ clinic: clinicData })
        });
    }

    async updateClinic(id, clinicData) {
        return this.request(API_CONFIG.ENDPOINTS.CLINIC(id), {
            method: 'PATCH',
            body: JSON.stringify({ clinic: clinicData })
        });
    }

    async deleteClinic(id) {
        return this.request(API_CONFIG.ENDPOINTS.CLINIC(id), {
            method: 'DELETE'
        });
    }

    // PATIENT ACCESS
    async createAccessRequest(patientId, clinicId, medicalDataId) {
        return this.request('/patient-access/request', {
            method: 'POST',
            body: JSON.stringify({
                patient_id: Number(patientId),
                clinic_id: Number(clinicId),
                medical_data_id: Number(medicalDataId)
            })
        });
    }

    async approveAccessRequest(requestId, employeeId) {
        return this.request(`/patient-access/approve/${requestId}`, {
            method: 'POST',
            body: JSON.stringify({
                employee_id: Number(employeeId)
            })
        });
    }

    async rejectAccessRequest(requestId, employeeId) {
        return this.request(`/patient-access/reject/${requestId}`, {
            method: 'POST',
            body: JSON.stringify({
                employee_id: Number(employeeId)
            })
        });
    }

    async listAccessRequests(filters = {}) {
        const params = new URLSearchParams(filters);
        return this.request(`/patient-access/requests?${params}`);
    }

    async getAccessRequest(id) {
        return this.request(`/patient-access/request/${id}`);
    }

    async getTemporaryData(token) {
        return this.request(`/patient-access/data?token=${encodeURIComponent(token)}`);
    }

    async revokeAccess(token) {
        return this.request('/patient-access/revoke', {
            method: 'POST',
            body: JSON.stringify({
                access_token: token
            })
        });
    }
}

const api = new APIClient();