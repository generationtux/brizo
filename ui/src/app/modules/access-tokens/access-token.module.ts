import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { RouterModule, Routes } from '@angular/router';
import { SharedModule } from '../shared.module';

import { AuthGuard } from '../auth/auth.guard';
import { AccessTokenComponent } from './access-token.component';

const accessTokensRoutes: Routes = [
  { path: 'access-tokens', component: AccessTokenComponent, canActivate: [AuthGuard] },
];

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    RouterModule.forChild(accessTokensRoutes),
    SharedModule,
  ],
  declarations: [
    AccessTokenComponent,
  ],
  exports: [
    AccessTokenComponent,
    RouterModule,
  ],
})

export class AccessTokenModule {};
