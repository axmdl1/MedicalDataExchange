window.APP_CONFIG = {
    BASE_URL: "http://localhost:8099", // your core-service API
    ENDPOINTS: {
        login: "/users/login",
        register: "/users",
        listMedicalData: (userId) => `/medical-data?user_id=${userId}`,
        deleteMedicalData: (id) => `/medical-data/${id}`,
    },
};
