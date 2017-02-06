import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { RouterModule, Routes } from '@angular/router';
import { SharedModule } from '../shared.module';

import { AuthGuard } from '../auth/auth.guard';
import { InviteComponent } from './invite.component';

const inviteRoutes: Routes = [
  { path: 'users/invite', component: InviteComponent, canActivate: [AuthGuard] },
];

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    RouterModule.forChild(inviteRoutes),
    SharedModule,
  ],
  declarations: [
    InviteComponent,
  ],
  exports: [
    InviteComponent,
    RouterModule,
  ],
})

export class InviteModule {};
