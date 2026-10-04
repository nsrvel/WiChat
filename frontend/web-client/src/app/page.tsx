import { ThemeToggle } from "@/components/theme-toggle";
import { Button } from "@/components/ui/button";

export default function HomePage() {
	return (
		<main className="flex flex-1 flex-col items-center justify-center gap-6 p-8">
			<div className="flex max-w-md flex-col items-center gap-2 text-center">
				<h1 className="text-3xl font-semibold tracking-tight">WiChat</h1>
				<p className="text-muted-foreground text-sm">
					Web client scaffold — auth and workspace slices come next.
				</p>
			</div>
			<div className="flex items-center gap-2">
				<Button disabled>Sign in (soon)</Button>
				<ThemeToggle />
			</div>
		</main>
	);
}
