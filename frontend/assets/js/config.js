// API Configuration
const API_CONFIG = {
    BASE_URL: 'http://localhost:8099',
    ENDPOINTS: {
        // Auth
        LOGIN: '/users/login',

        // Users
        USERS: '/users',
        USER: (id) => `/users/${id}`,

        // Medical Data
        MEDICAL_DATA: '/medical-data',
        MEDICAL_DATA_ITEM: (id) => `/medical-data/${id}`,

        // Transfers
        TRANSFERS: '/transfers',
        TRANSFER: (id) => `/transfers/${id}`,
        TRANSFER_DECISION: '/transfers/decision',

        // Clinics
        CLINICS: '/clinics',
        CLINIC: (id) => `/clinics/${id}`
    }
};

// User roles
const ROLES = {
    PATIENT: 'patient',
    EMPLOYEE: 'employee',
    ADMIN: 'admin'
};

// Transfer statuses
const TRANSFER_STATUS = {
    PENDING: 'pending',
    CONFIRMED: 'confirmed',
    REJECTED: 'rejected'
};
