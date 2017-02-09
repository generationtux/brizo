import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { RouterModule, Routes } from '@angular/router';
import { HttpModule } from '@angular/http';

import { SharedModule } from './modules/shared.module';

import { ApplicationModule } from './modules/applications/application.module';
import { AuthModule } from './modules/auth/auth.module';
import { DashboardModule } from './modules/dashboard/dashboard.module';
import { EnvironmentModule } from './modules/environments/environment.module';
import { HelpModule } from './modules/help/help.module';
import { InviteModule } from './modules/invite/invite.module';
import { LoginModule } from './modules/login/login.module';
import { VersionModule } from './modules/versions/version.module';

import { AppComponent } from './app.component';
import { AuthGuard } from './modules/auth/auth.guard';
import { AuthService } from './modules/auth/auth.service';

// Routes are delegated to each module
const routes: Routes = [];

@NgModule({
    imports: [
      ApplicationModule,
      AuthModule,
      BrowserModule,
      DashboardModule,
      EnvironmentModule,
      FormsModule,
      HelpModule,
      HttpModule,
      InviteModule,
      LoginModule,
      ReactiveFormsModule,
      RouterModule.forRoot(routes),
      SharedModule,
      VersionModule,
    ],
    providers:    [ AuthGuard ],
    declarations: [ AppComponent],
    bootstrap:    [ AppComponent ]
})

export class AppModule {};
