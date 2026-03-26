import { BusinessList } from '@/services/auth/business/ui';

export default function BusinessesPageRoute() {
    return (
        <div className="p-6 space-y-6">
            <BusinessList />
        </div>
    );
}
