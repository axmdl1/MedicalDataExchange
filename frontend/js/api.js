const cfg = window.APP_CONFIG;

// helper to attach JWT
function getAuthHeaders() {
    const token = localStorage.getItem("token");
    return token ? { "Authorization": `Bearer ${token}` } : {};
}

// LOGIN
export async function apiLogin(email, password) {
    const res = await fetch(cfg.BASE_URL + cfg.ENDPOINTS.login, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
    });
    if (!res.ok) throw new Error(await res.text());
    const data = await res.json();

    // ✅ save correct token
    const token = data.access_token || data.token;
    if (token) localStorage.setItem("token", token);

    if (data.user) localStorage.setItem("user", JSON.stringify(data.user));
    return data;
}

// REGISTER
export async function apiRegister(user) {
    const res = await fetch(cfg.BASE_URL + cfg.ENDPOINTS.register, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ user }),
    });
    if (!res.ok) throw new Error(await res.text());
    return res.json();
}

// MEDICAL DATA
export async function apiListMedicalData(userId) {
    const res = await fetch(cfg.BASE_URL + cfg.ENDPOINTS.listMedicalData(userId), {
        headers: { "Content-Type": "application/json", ...getAuthHeaders() },
    });
    if (!res.ok) throw new Error(await res.text());
    return res.json();
}

export async function apiDeleteMedicalData(id) {
    const res = await fetch(cfg.BASE_URL + cfg.ENDPOINTS.deleteMedicalData(id), {
        method: "DELETE",
        headers: { ...getAuthHeaders() },
    });
    if (!res.ok) throw new Error(await res.text());
    return res.json();
}
