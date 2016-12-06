/**
 * Angular Bootstrapping
 */
import { platformBrowserDynamic } from '@angular/platform-browser-dynamic';

/**
 * App Module
 * ==========
 * Import our tope level module
 */
import { AppModule } from './app/app.module'

platformBrowserDynamic().bootstrapModule(AppModule);
