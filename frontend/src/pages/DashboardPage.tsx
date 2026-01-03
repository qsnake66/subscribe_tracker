import { useEffect, useMemo, useState, type FormEvent } from 'react';
import { CalendarClock, CreditCard, LogOut, Pencil, Plus, RefreshCcw, Trash2 } from 'lucide-react';
import { cn } from '../lib/utils';
import {
  createSubscription,
  deleteSubscription,
  getSubscriptions,
  updateSubscription,
  type Subscription,
} from '../lib/api';
import { clearAuth, getAuthUser } from '../lib/auth';

const emptyForm: Omit<Subscription, 'id'> = {
  service_name: '',
  bank_name: '',
  card_last4: '',
  billing_cycle: 'monthly',
  charge_date: '',
};

export default function DashboardPage() {
  const user = getAuthUser();
  const [subscriptions, setSubscriptions] = useState<Subscription[]>([]);
  const [form, setForm] = useState<Omit<Subscription, 'id'>>(emptyForm);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [saving, setSaving] = useState(false);

  const totalCount = subscriptions.length;

  const upcoming = useMemo(() => {
    return subscriptions.slice(0, 3);
  }, [subscriptions]);

  useEffect(() => {
    void loadSubscriptions();
  }, []);

  const loadSubscriptions = async () => {
    setLoading(true);
    setError('');
    try {
      const data = await getSubscriptions();
      setSubscriptions(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Не удалось загрузить подписки');
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (field: keyof Omit<Subscription, 'id'>, value: string) => {
    setForm((prev) => ({ ...prev, [field]: value }));
  };

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault();
    setSaving(true);
    setError('');

    try {
      if (editingId) {
        await updateSubscription(editingId, form);
      } else {
        await createSubscription(form);
      }
      setForm(emptyForm);
      setEditingId(null);
      await loadSubscriptions();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка сохранения');
    } finally {
      setSaving(false);
    }
  };

  const handleEdit = (item: Subscription) => {
    setEditingId(item.id);
    setForm({
      service_name: item.service_name,
      bank_name: item.bank_name,
      card_last4: item.card_last4,
      billing_cycle: item.billing_cycle,
      charge_date: item.charge_date,
    });
  };

  const handleDelete = async (id: string) => {
    setError('');
    try {
      await deleteSubscription(id);
      await loadSubscriptions();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка удаления');
    }
  };

  const handleLogout = () => {
    clearAuth();
    window.location.href = '/';
  };

  return (
    <div className="min-h-screen bg-background text-white relative overflow-hidden">
      <div className="absolute -top-20 left-[-10%] h-[420px] w-[420px] rounded-full bg-secondary/20 blur-[120px]" />
      <div className="absolute bottom-[-20%] right-[-10%] h-[480px] w-[480px] rounded-full bg-accent/20 blur-[140px]" />

      <div className="relative z-10 mx-auto max-w-6xl px-4 py-10 sm:px-6 lg:px-8">
        <header className="flex flex-col gap-6 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <p className="text-sm uppercase tracking-[0.25em] text-gray-400">Личный кабинет</p>
            <h1 className="mt-2 text-3xl font-bold">Подписки {user?.name ? `, ${user.name}` : ''}</h1>
            <p className="mt-2 text-sm text-gray-400">Управляй сервисами и датами списаний в одном месте.</p>
          </div>
          <div className="flex gap-3">
            <button
              onClick={() => void loadSubscriptions()}
              className="inline-flex items-center gap-2 rounded-lg border border-white/10 px-4 py-2 text-sm text-gray-200 hover:border-white/30 hover:text-white transition"
            >
              <RefreshCcw className="h-4 w-4" />
              Обновить
            </button>
            <button
              onClick={handleLogout}
              className="inline-flex items-center gap-2 rounded-lg bg-white/5 px-4 py-2 text-sm text-gray-200 hover:bg-white/10 transition"
            >
              <LogOut className="h-4 w-4" />
              Выйти
            </button>
          </div>
        </header>

        <section className="mt-10 grid gap-6 lg:grid-cols-[2fr_1fr]">
          <div className="space-y-6">
            <div className="grid gap-4 sm:grid-cols-3">
              <div className="rounded-2xl border border-white/10 bg-surface/70 p-4 shadow-xl">
                <p className="text-xs uppercase tracking-[0.3em] text-gray-400">Всего</p>
                <p className="mt-3 text-2xl font-semibold">{totalCount}</p>
              </div>
              <div className="rounded-2xl border border-white/10 bg-surface/70 p-4 shadow-xl">
                <p className="text-xs uppercase tracking-[0.3em] text-gray-400">Ближайшие</p>
                <p className="mt-3 text-2xl font-semibold">{upcoming.length}</p>
              </div>
              <div className="rounded-2xl border border-white/10 bg-surface/70 p-4 shadow-xl">
                <p className="text-xs uppercase tracking-[0.3em] text-gray-400">Статус</p>
                <p className="mt-3 text-sm text-primary">Контроль расходов</p>
              </div>
            </div>

            <div className="rounded-2xl border border-white/10 bg-surface/60 p-6 shadow-2xl">
              <div className="flex items-center justify-between">
                <h2 className="text-xl font-semibold">Список подписок</h2>
                <span className="text-sm text-gray-400">{loading ? 'Загрузка...' : `${subscriptions.length} записей`}</span>
              </div>

              {error && (
                <div className="mt-4 rounded-lg border border-red-500/40 bg-red-500/10 px-4 py-3 text-sm text-red-200">
                  {error}
                </div>
              )}

              <div className="mt-6 space-y-4">
                {!loading && subscriptions.length === 0 && (
                  <div className="rounded-xl border border-dashed border-white/10 px-6 py-10 text-center text-sm text-gray-400">
                    Пока нет подписок. Добавь первую справа.
                  </div>
                )}

                {subscriptions.map((item) => (
                  <div
                    key={item.id}
                    className="flex flex-col gap-4 rounded-xl border border-white/10 bg-background/40 p-4 sm:flex-row sm:items-center sm:justify-between"
                  >
                    <div>
                      <p className="text-lg font-semibold">{item.service_name}</p>
                      <p className="mt-1 text-sm text-gray-400">
                        {item.bank_name} · **** {item.card_last4}
                      </p>
                      <div className="mt-3 flex flex-wrap items-center gap-3 text-xs text-gray-300">
                        <span className="inline-flex items-center gap-2 rounded-full border border-white/10 px-3 py-1">
                          <CalendarClock className="h-3.5 w-3.5" />
                          {item.charge_date}
                        </span>
                        <span
                          className={cn(
                            'inline-flex items-center gap-2 rounded-full border border-white/10 px-3 py-1 capitalize',
                            item.billing_cycle === 'monthly' ? 'text-primary' : 'text-secondary'
                          )}
                        >
                          <CreditCard className="h-3.5 w-3.5" />
                          {item.billing_cycle === 'monthly' ? 'ежемесячно' : 'ежегодно'}
                        </span>
                      </div>
                    </div>
                    <div className="flex items-center gap-2">
                      <button
                        onClick={() => handleEdit(item)}
                        className="inline-flex items-center gap-2 rounded-lg border border-white/10 px-3 py-2 text-xs text-gray-200 hover:border-white/30"
                      >
                        <Pencil className="h-4 w-4" />
                        Изменить
                      </button>
                      <button
                        onClick={() => void handleDelete(item.id)}
                        className="inline-flex items-center gap-2 rounded-lg border border-red-500/40 px-3 py-2 text-xs text-red-200 hover:border-red-500"
                      >
                        <Trash2 className="h-4 w-4" />
                        Удалить
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>

          <aside className="rounded-2xl border border-white/10 bg-surface/70 p-6 shadow-2xl">
            <div className="flex items-center justify-between">
              <h2 className="text-lg font-semibold">{editingId ? 'Редактирование' : 'Новая подписка'}</h2>
              <span className="inline-flex items-center gap-2 rounded-full border border-white/10 px-3 py-1 text-xs text-gray-300">
                <Plus className="h-3.5 w-3.5" />
                MVP
              </span>
            </div>

            <form className="mt-6 space-y-4" onSubmit={handleSubmit}>
              <div>
                <label className="text-xs uppercase tracking-[0.2em] text-gray-400">Сервис</label>
                <input
                  value={form.service_name}
                  onChange={(e) => handleChange('service_name', e.target.value)}
                  className="mt-2 w-full rounded-lg border border-white/10 bg-background/40 px-3 py-2 text-sm text-white placeholder-gray-500 focus:border-primary focus:outline-none"
                  placeholder="Netflix, Spotify"
                  required
                />
              </div>
              <div>
                <label className="text-xs uppercase tracking-[0.2em] text-gray-400">Банк</label>
                <input
                  value={form.bank_name}
                  onChange={(e) => handleChange('bank_name', e.target.value)}
                  className="mt-2 w-full rounded-lg border border-white/10 bg-background/40 px-3 py-2 text-sm text-white placeholder-gray-500 focus:border-primary focus:outline-none"
                  placeholder="Tinkoff, Sber"
                  required
                />
              </div>
              <div>
                <label className="text-xs uppercase tracking-[0.2em] text-gray-400">4 цифры карты</label>
                <input
                  value={form.card_last4}
                  onChange={(e) => handleChange('card_last4', e.target.value.replace(/\D/g, '').slice(0, 4))}
                  className="mt-2 w-full rounded-lg border border-white/10 bg-background/40 px-3 py-2 text-sm text-white placeholder-gray-500 focus:border-primary focus:outline-none"
                  placeholder="1234"
                  required
                />
              </div>
              <div>
                <label className="text-xs uppercase tracking-[0.2em] text-gray-400">Дата списания</label>
                <input
                  type="date"
                  value={form.charge_date}
                  onChange={(e) => handleChange('charge_date', e.target.value)}
                  className="mt-2 w-full rounded-lg border border-white/10 bg-background/40 px-3 py-2 text-sm text-white placeholder-gray-500 focus:border-primary focus:outline-none"
                  required
                />
              </div>
              <div>
                <label className="text-xs uppercase tracking-[0.2em] text-gray-400">Периодичность</label>
                <div className="mt-2 flex gap-2">
                  {(['monthly', 'yearly'] as const).map((value) => (
                    <button
                      key={value}
                      type="button"
                      onClick={() => handleChange('billing_cycle', value)}
                      className={cn(
                        'flex-1 rounded-lg border px-3 py-2 text-xs uppercase tracking-[0.2em] transition',
                        form.billing_cycle === value
                          ? 'border-primary bg-primary/10 text-primary'
                          : 'border-white/10 text-gray-400 hover:border-white/30'
                      )}
                    >
                      {value === 'monthly' ? 'ежемесячно' : 'ежегодно'}
                    </button>
                  ))}
                </div>
              </div>

              <div className="flex gap-2">
                <button
                  type="submit"
                  disabled={saving}
                  className="flex-1 rounded-lg bg-primary px-4 py-3 text-sm font-semibold text-white hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-60"
                >
                  {saving ? 'Сохранение...' : editingId ? 'Сохранить' : 'Добавить'}
                </button>
                {editingId && (
                  <button
                    type="button"
                    onClick={() => {
                      setEditingId(null);
                      setForm(emptyForm);
                    }}
                    className="rounded-lg border border-white/10 px-4 py-3 text-sm text-gray-300 hover:border-white/30"
                  >
                    Отмена
                  </button>
                )}
              </div>
            </form>
          </aside>
        </section>
      </div>
    </div>
  );
}
