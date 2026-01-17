#!/bin/bash

# Script para asociar el IAM Role correcto a la instancia EC2
# Probability - Asociar probability-ec2-ecr-pull-profile

INSTANCE_ID="i-0f3284d2a87127e57"
INSTANCE_PROFILE_NAME="probability-ec2-ecr-pull-profile"

echo "🔧 Asociando IAM Role a la instancia EC2..."

# Primero, desasociar el rol actual si existe
echo "📋 Verificando rol actual..."
CURRENT_PROFILE=$(aws ec2 describe-instances \
  --instance-ids $INSTANCE_ID \
  --query 'Reservations[0].Instances[0].IamInstanceProfile.Arn' \
  --output text 2>/dev/null)

if [ "$CURRENT_PROFILE" != "None" ] && [ -n "$CURRENT_PROFILE" ]; then
  echo "⚠️  Instancia tiene un rol asociado: $CURRENT_PROFILE"
  echo "🔄 Desasociando rol actual..."
  
  ASSOCIATION_ID=$(aws ec2 describe-iam-instance-profile-associations \
    --filters "Name=instance-id,Values=$INSTANCE_ID" \
    --query 'IamInstanceProfileAssociations[0].AssociationId' \
    --output text 2>/dev/null)
  
  if [ "$ASSOCIATION_ID" != "None" ] && [ -n "$ASSOCIATION_ID" ]; then
    aws ec2 disassociate-iam-instance-profile \
      --association-id $ASSOCIATION_ID
    echo "✅ Rol anterior desasociado"
    sleep 2
  fi
fi

# Asociar el nuevo rol
echo "🔗 Asociando nuevo rol: $INSTANCE_PROFILE_NAME"
aws ec2 associate-iam-instance-profile \
  --instance-id $INSTANCE_ID \
  --iam-instance-profile Name=$INSTANCE_PROFILE_NAME

if [ $? -eq 0 ]; then
  echo "✅ Rol asociado correctamente!"
  echo ""
  echo "📝 Verificando asociación..."
  aws ec2 describe-iam-instance-profile-associations \
    --filters "Name=instance-id,Values=$INSTANCE_ID" \
    --query 'IamInstanceProfileAssociations[0].[IamInstanceProfile.Arn,State]' \
    --output table
else
  echo "❌ Error al asociar el rol"
  exit 1
fi


