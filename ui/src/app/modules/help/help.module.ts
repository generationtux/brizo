import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Routes } from '@angular/router';
import { SharedModule } from '../shared.module';

import { AuthGuard } from '../auth/auth.guard';
import { HelpComponent } from './help.component';

const helpRoutes: Routes = [
  { path: 'help', component: HelpComponent, canActivate: [AuthGuard] },
];

@NgModule({
  imports: [
    CommonModule,
    RouterModule.forChild(helpRoutes),
    SharedModule,
  ],
  declarations: [
    HelpComponent,
  ],
  exports: [
    HelpComponent,
    RouterModule,
  ],
})

export class HelpModule {}
