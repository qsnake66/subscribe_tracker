import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { Link, useNavigate } from 'react-router-dom';
import { cn } from '../lib/utils';
import { Lock, Mail, User, ArrowRight } from 'lucide-react';
import { registerUser } from '../lib/api';

const registerSchema = z.object({
    name: z.string().min(2, 'Name must be at least 2 characters'),
    email: z.string().email('Please enter a valid email address'),
    password: z
        .string()
        .min(8, 'Password must be at least 8 characters')
        .regex(/[A-Z]/, 'Password must contain at least one uppercase letter')
        .regex(/[a-z]/, 'Password must contain at least one lowercase letter')
        .regex(/[0-9]/, 'Password must contain at least one number'),
    confirmPassword: z.string(),
}).refine((data) => data.password === data.confirmPassword, {
    message: "Passwords don't match",
    path: ["confirmPassword"],
});

type RegisterFormData = z.infer<typeof registerSchema>;

export default function RegisterPage() {
    const navigate = useNavigate();
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);
    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm<RegisterFormData>({
        resolver: zodResolver(registerSchema),
    });

    const onSubmit = async (data: RegisterFormData) => {
        setError('');
        setLoading(true);
        try {
            await registerUser(data.name, data.email, data.password);
            navigate('/app');
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Не удалось создать аккаунт');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="flex min-h-screen items-center justify-center bg-background px-4 py-12 sm:px-6 lg:px-8 relative overflow-hidden">
            <div className="absolute top-[-10%] right-[-10%] h-[500px] w-[500px] rounded-full bg-primary/10 blur-[100px]" />
            <div className="absolute bottom-[-10%] left-[-10%] h-[500px] w-[500px] rounded-full bg-accent/20 blur-[100px]" />

            <div className="w-full max-w-md space-y-8 relative z-10 bg-surface/50 backdrop-blur-xl p-8 rounded-2xl border border-white/10 shadow-2xl">
                <div className="text-center">
                    <h2 className="mt-2 text-3xl font-bold tracking-tight text-white">
                        Create Account
                    </h2>
                    <p className="mt-2 text-sm text-gray-400">
                        Start tracking all your subscriptions in one place
                    </p>
                </div>

                {error && (
                    <div className="rounded-lg border border-red-500/40 bg-red-500/10 px-4 py-3 text-sm text-red-200">
                        {error}
                    </div>
                )}

                <form className="mt-8 space-y-6" onSubmit={handleSubmit(onSubmit)}>
                    <div className="space-y-4">
                        <div>
                            <label htmlFor="name" className="sr-only">
                                Full Name
                            </label>
                            <div className="relative">
                                <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                                    <User className="h-5 w-5 text-gray-500" aria-hidden="true" />
                                </div>
                                <input
                                    id="name"
                                    type="text"
                                    autoComplete="name"
                                    className={cn(
                                        "block w-full rounded-lg border border-white/10 bg-background/50 py-3 pl-10 pr-3 text-white placeholder-gray-500 focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary sm:text-sm transition-all duration-200",
                                        errors.name && "border-red-500 focus:border-red-500 focus:ring-red-500"
                                    )}
                                    placeholder="Full Name"
                                    {...register('name')}
                                />
                            </div>
                            {errors.name && (
                                <p className="mt-1 text-xs text-red-500 pl-1">{errors.name.message}</p>
                            )}
                        </div>

                        <div>
                            <label htmlFor="email" className="sr-only">
                                Email address
                            </label>
                            <div className="relative">
                                <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                                    <Mail className="h-5 w-5 text-gray-500" aria-hidden="true" />
                                </div>
                                <input
                                    id="email"
                                    type="email"
                                    autoComplete="email"
                                    className={cn(
                                        "block w-full rounded-lg border border-white/10 bg-background/50 py-3 pl-10 pr-3 text-white placeholder-gray-500 focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary sm:text-sm transition-all duration-200",
                                        errors.email && "border-red-500 focus:border-red-500 focus:ring-red-500"
                                    )}
                                    placeholder="Email address"
                                    {...register('email')}
                                />
                            </div>
                            {errors.email && (
                                <p className="mt-1 text-xs text-red-500 pl-1">{errors.email.message}</p>
                            )}
                        </div>

                        <div>
                            <label htmlFor="password" className="sr-only">
                                Password
                            </label>
                            <div className="relative">
                                <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                                    <Lock className="h-5 w-5 text-gray-500" aria-hidden="true" />
                                </div>
                                <input
                                    id="password"
                                    type="password"
                                    autoComplete="new-password"
                                    className={cn(
                                        "block w-full rounded-lg border border-white/10 bg-background/50 py-3 pl-10 pr-3 text-white placeholder-gray-500 focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary sm:text-sm transition-all duration-200",
                                        errors.password && "border-red-500 focus:border-red-500 focus:ring-red-500"
                                    )}
                                    placeholder="Password"
                                    {...register('password')}
                                />
                            </div>
                            {errors.password && (
                                <p className="mt-1 text-xs text-red-500 pl-1">{errors.password.message}</p>
                            )}
                        </div>

                        <div>
                            <label htmlFor="confirmPassword" className="sr-only">
                                Confirm Password
                            </label>
                            <div className="relative">
                                <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                                    <Lock className="h-5 w-5 text-gray-500" aria-hidden="true" />
                                </div>
                                <input
                                    id="confirmPassword"
                                    type="password"
                                    autoComplete="new-password"
                                    className={cn(
                                        "block w-full rounded-lg border border-white/10 bg-background/50 py-3 pl-10 pr-3 text-white placeholder-gray-500 focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary sm:text-sm transition-all duration-200",
                                        errors.confirmPassword && "border-red-500 focus:border-red-500 focus:ring-red-500"
                                    )}
                                    placeholder="Confirm Password"
                                    {...register('confirmPassword')}
                                />
                            </div>
                            {errors.confirmPassword && (
                                <p className="mt-1 text-xs text-red-500 pl-1">{errors.confirmPassword.message}</p>
                            )}
                        </div>
                    </div>

                    <div>
                        <button
                            type="submit"
                            disabled={loading}
                            className="group relative flex w-full justify-center rounded-lg bg-primary py-3 px-4 text-sm font-semibold text-white hover:bg-primary/90 focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2 focus:ring-offset-gray-900 transition-all duration-200 disabled:cursor-not-allowed disabled:opacity-60"
                        >
                            {loading ? 'Creating...' : 'create account'}
                            <span className="absolute inset-y-0 right-0 flex items-center pr-3">
                                <ArrowRight className="h-4 w-4 text-white/50 group-hover:text-white transition-colors" />
                            </span>
                        </button>
                    </div>
                </form>

                <div className="text-center text-sm">
                    <p className="text-gray-400">
                        Already have an account?{' '}
                        <Link to="/" className="font-semibold text-primary hover:text-primary/80 transition-colors">
                            Sign in
                        </Link>
                    </p>
                </div>
            </div>
        </div>
    );
}
