import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { AuthComponent } from './auth.component';
import { AuthGuard } from './auth.guard';
import { AuthService } from './auth.service';

const authRoutes: Routes = [
  { path: 'auth', component: AuthComponent },
];

@NgModule({
  imports: [
    RouterModule.forChild(authRoutes),
  ],
  declarations: [
    AuthComponent,
  ],
  providers: [
    AuthGuard,
    AuthService,
  ],
  exports: [
    AuthComponent,
    RouterModule,
  ],
})

export class AuthModule {};
