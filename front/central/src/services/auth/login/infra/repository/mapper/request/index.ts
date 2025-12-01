// Interfaces
export interface LoginRequest {
    email: string;
    password: string;
}

export interface ChangePasswordRequest {
    current_password: string;
    new_password: string;
}

export interface GeneratePasswordRequest {
    user_id?: number;
}

export interface GenerateBusinessTokenRequest {
    business_id: number;
}

// Mappers
export const mapLoginRequest = (req: LoginRequest): any => {
    return {
        email: req.email,
        password: req.password,
    };
};

export const mapChangePasswordRequest = (req: ChangePasswordRequest): any => {
    return {
        current_password: req.current_password,
        new_password: req.new_password,
    };
};

export const mapGeneratePasswordRequest = (req: GeneratePasswordRequest): any => {
    return {
        user_id: req.user_id,
    };
};

export const mapGenerateBusinessTokenRequest = (req: GenerateBusinessTokenRequest): any => {
    return {
        business_id: req.business_id,
    };
};
