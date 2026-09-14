import { Suspense } from "react";
import { NotificationDashboard } from "@/services/modules/notification-config/ui/components/NotificationDashboard";

export default function NotificationConfigPage() {
  return (
    <Suspense fallback={null}>
      <NotificationDashboard />
    </Suspense>
  );
}
